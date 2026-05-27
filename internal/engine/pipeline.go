package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/internal/services"
	logsyslog "github.com/logcat/logcat/internal/syslog"
)

// =============================================================================
// Pipeline item that flows through the processing stages.
// =============================================================================

// PipelineItem carries a log entry through every stage of the pipeline.
type PipelineItem struct {
	ParsedLog    *logsyslog.ParsedLog
	Device       *models.Device
	ParsedData   map[string]interface{}
	ParseOK      bool
	ParseErr     string
	FilterResult *services.FilterResult
	Deduped      bool
	HighFreq     bool
	Trace        *models.LogTraceInfo
}

// =============================================================================
// Pipeline metrics
// =============================================================================

// PipelineMetrics exposes counters and queue depths for observability.
type PipelineMetrics struct {
	// Queue depths (approximate, via len())
	RawQueueDepth    int `json:"rawQueueDepth"`
	ParsedQueueDepth int `json:"parsedQueueDepth"`
	DBQueueDepth     int `json:"dbQueueDepth"`
	PushQueueDepth   int `json:"pushQueueDepth"`

	// Stage counts
	ParseProcessed  int64 `json:"parseProcessed"`
	ParseErrors     int64 `json:"parseErrors"`
	FilterProcessed int64 `json:"filterProcessed"`
	FilterDropped   int64 `json:"filterDropped"`
	DBWritten       int64 `json:"dbWritten"`
	DBErrors        int64 `json:"dbErrors"`
	PushProcessed   int64 `json:"pushProcessed"`
	PushErrors      int64 `json:"pushErrors"`

	// Drop counts due to queue full
	RawDropped  int64 `json:"rawDropped"`
	DBDropped   int64 `json:"dbDropped"`
	PushDropped int64 `json:"pushDropped"`
}

// =============================================================================
// Pipeline
// =============================================================================

// Pipeline implements the complete async log processing pipeline.
type Pipeline struct {
	// Channel stages
	rawCh    chan *logsyslog.ParsedLog // input from syslog receiver
	parsedCh chan *PipelineItem        // after parse workers
	dbCh     chan *PipelineItem        // for DB writer
	pushCh   chan *PipelineItem        // for push workers

	// Worker counts
	parseWorkers  int
	filterWorkers int
	pushWorkers   int

	// Lifecycle
	stopCh  chan struct{}
	wg      sync.WaitGroup
	running bool
	mu      sync.Mutex

	// Metrics
	metrics PipelineMetrics

	// Services
	deviceService        *services.DeviceService
	parseService         *services.ParseService
	filterService        *services.FilterService
	dedupService         *services.DedupService
	aggregateService     *services.AggregateService
	highFreqService      *services.HighFreqService
	desensitizeService   *services.DesensitizeService
	pushService          *services.PushService
	emailService         *services.EmailService
	syslogForwardService *services.SyslogForwardService
	alertService         *services.AlertService
	traceService         *services.TraceService
}

// =============================================================================
// PipelineServices bundles all services needed by the pipeline.
// =============================================================================

// PipelineServices holds references to every service the pipeline uses.
type PipelineServices struct {
	DeviceService        *services.DeviceService
	ParseService         *services.ParseService
	FilterService        *services.FilterService
	DedupService         *services.DedupService
	AggregateService     *services.AggregateService
	HighFreqService      *services.HighFreqService
	DesensitizeService   *services.DesensitizeService
	PushService          *services.PushService
	EmailService         *services.EmailService
	SyslogForwardService *services.SyslogForwardService
	AlertService         *services.AlertService
	TraceService         *services.TraceService
}

// NewPipeline creates a new Pipeline.
func NewPipeline(cfg *config.Config, svc *PipelineServices) *Pipeline {
	cap := cfg.Queue.Capacity
	if cap <= 0 {
		cap = 10000
	}

	return &Pipeline{
		rawCh:    make(chan *logsyslog.ParsedLog, cap),
		parsedCh: make(chan *PipelineItem, cap),
		dbCh:     make(chan *PipelineItem, cap),
		pushCh:   make(chan *PipelineItem, cap),

		parseWorkers:  cfg.Worker.ParseWorkers,
		filterWorkers: cfg.Worker.FilterWorkers,
		pushWorkers:   cfg.Worker.PushWorkers,

		stopCh: make(chan struct{}),

		deviceService:        svc.DeviceService,
		parseService:         svc.ParseService,
		filterService:        svc.FilterService,
		dedupService:         svc.DedupService,
		aggregateService:     svc.AggregateService,
		highFreqService:      svc.HighFreqService,
		desensitizeService:   svc.DesensitizeService,
		pushService:          svc.PushService,
		emailService:         svc.EmailService,
		syslogForwardService: svc.SyslogForwardService,
		alertService:         svc.AlertService,
		traceService:         svc.TraceService,
	}
}

// RawChannel returns the input channel that the syslog receiver writes to.
func (p *Pipeline) RawChannel() chan<- *logsyslog.ParsedLog {
	return p.rawCh
}

// Start initializes and starts all pipeline worker goroutines.
func (p *Pipeline) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return
	}

	p.running = true
	log.Printf("[pipeline] Starting with %d parse, %d filter, %d push workers",
		p.parseWorkers, p.filterWorkers, p.pushWorkers)

	// Parse workers
	for i := 0; i < p.parseWorkers; i++ {
		p.wg.Add(1)
		go p.parseWorker(i)
	}

	// Filter workers
	for i := 0; i < p.filterWorkers; i++ {
		p.wg.Add(1)
		go p.filterWorker(i)
	}

	// DB writer (single goroutine, batch insert)
	p.wg.Add(1)
	go p.dbWriter()

	// DB batch flush ticker goroutine
	p.wg.Add(1)
	go p.dbFlushTicker()

	// Push workers
	for i := 0; i < p.pushWorkers; i++ {
		p.wg.Add(1)
		go p.pushWorker(i)
	}
}

// Stop gracefully shuts down the pipeline: stops accepting new input,
// drains queues, and waits for workers to finish.
func (p *Pipeline) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return
	}

	p.running = false

	// Signal all workers to stop
	close(p.stopCh)

	// Wait for all workers to drain and exit
	p.wg.Wait()

	// Flush any remaining batch
	p.flushBatch()

	log.Println("[pipeline] Stopped")
}

// IsRunning returns whether the pipeline is running.
func (p *Pipeline) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

// Metrics returns a snapshot of pipeline operational metrics.
func (p *Pipeline) Metrics() PipelineMetrics {
	return PipelineMetrics{
		RawQueueDepth:    len(p.rawCh),
		ParsedQueueDepth: len(p.parsedCh),
		DBQueueDepth:     len(p.dbCh),
		PushQueueDepth:   len(p.pushCh),

		ParseProcessed:  atomic.LoadInt64(&p.metrics.ParseProcessed),
		ParseErrors:     atomic.LoadInt64(&p.metrics.ParseErrors),
		FilterProcessed: atomic.LoadInt64(&p.metrics.FilterProcessed),
		FilterDropped:   atomic.LoadInt64(&p.metrics.FilterDropped),
		DBWritten:       atomic.LoadInt64(&p.metrics.DBWritten),
		DBErrors:        atomic.LoadInt64(&p.metrics.DBErrors),
		PushProcessed:   atomic.LoadInt64(&p.metrics.PushProcessed),
		PushErrors:      atomic.LoadInt64(&p.metrics.PushErrors),

		RawDropped:  atomic.LoadInt64(&p.metrics.RawDropped),
		DBDropped:   atomic.LoadInt64(&p.metrics.DBDropped),
		PushDropped: atomic.LoadInt64(&p.metrics.PushDropped),
	}
}

// =============================================================================
// Stage 1: Parse workers
// =============================================================================

func (p *Pipeline) parseWorker(id int) {
	defer p.wg.Done()
	defer p.recoverPanic("parseWorker", id)

	log.Printf("[pipeline] Parse worker %d started", id)

	for {
		select {
		case <-p.stopCh:
			// Drain remaining items in rawCh before exiting
			p.drainParse(id)
			log.Printf("[pipeline] Parse worker %d stopped", id)
			return
		case parsed := <-p.rawCh:
			p.processParse(id, parsed)
		}
	}
}

func (p *Pipeline) drainParse(id int) {
	for {
		select {
		case parsed := <-p.rawCh:
			p.processParse(id, parsed)
		default:
			return
		}
	}
}

func (p *Pipeline) processParse(id int, parsed *logsyslog.ParsedLog) {
	atomic.AddInt64(&p.metrics.ParseProcessed, 1)

	item := &PipelineItem{
		ParsedLog: parsed,
	}

	logEntry := parsed.SyslogLog

	// --- Device identification ---
	device, err := p.deviceService.GetDeviceByIP(logEntry.SourceIP)
	if err != nil {
		// Unknown device is OK; we still process the log
		item.Device = nil
	} else {
		item.Device = device
		item.ParsedLog.SyslogLog.DeviceID = &device.ID
		item.ParsedLog.SyslogLog.DeviceName = device.Name
	}

	// --- Trace creation ---
	trace, err := p.traceService.GetOrCreateTrace(logEntry.LogID)
	if err != nil {
		log.Printf("[pipeline] Trace creation failed for %s: %v", logEntry.LogID, err)
	}
	item.Trace = trace

	// --- Parse template lookup & field parsing ---
	if device != nil && device.ParseTemplateID != nil {
		result, err := p.parseService.ParseByTemplateID(*device.ParseTemplateID, logEntry.RawMessage)
		if err != nil {
			item.ParseOK = false
			item.ParseErr = err.Error()
			atomic.AddInt64(&p.metrics.ParseErrors, 1)

			if trace != nil {
				trace.ParseStatus = "error"
				trace.ParseError = err.Error()
				_ = p.traceService.UpdateTrace(logEntry.LogID, map[string]interface{}{
					"parse_status":      "error",
					"parse_error":       err.Error(),
					"parse_template_id": *device.ParseTemplateID,
				})
			}
		} else {
			item.ParseOK = result.Success
			item.ParsedData = result.Fields
			if !result.Success {
				item.ParseErr = result.Error
				atomic.AddInt64(&p.metrics.ParseErrors, 1)
			}

			if trace != nil {
				trace.ParseStatus = "done"
				if result.Success {
					b, _ := json.Marshal(result.Fields)
					trace.ParseResult = string(b)
				} else {
					trace.ParseError = result.Error
				}
				trace.ParseTemplateID = device.ParseTemplateID
				_ = p.traceService.UpdateTrace(logEntry.LogID, map[string]interface{}{
					"parse_status":      trace.ParseStatus,
					"parse_result":      trace.ParseResult,
					"parse_error":       trace.ParseError,
					"parse_template_id": device.ParseTemplateID,
				})
			}
		}
	} else {
		// No parse template: treat raw message as data
		item.ParsedData = map[string]interface{}{
			"raw_message": logEntry.RawMessage,
			"source_ip":   logEntry.SourceIP,
			"severity":    logEntry.Severity,
			"facility":    logEntry.Facility,
		}
		item.ParseOK = true
	}

	// Forward to filter stage
	sendOrDrop(p.parsedCh, item, &p.metrics.FilterDropped, "filter")
}

// =============================================================================
// Stage 2: Filter workers
// =============================================================================

func (p *Pipeline) filterWorker(id int) {
	defer p.wg.Done()
	defer p.recoverPanic("filterWorker", id)

	log.Printf("[pipeline] Filter worker %d started", id)

	for {
		select {
		case <-p.stopCh:
			p.drainFilter(id)
			log.Printf("[pipeline] Filter worker %d stopped", id)
			return
		case item := <-p.parsedCh:
			p.processFilter(item)
		}
	}
}

func (p *Pipeline) drainFilter(id int) {
	for {
		select {
		case item := <-p.parsedCh:
			p.processFilter(item)
		default:
			return
		}
	}
}

func (p *Pipeline) processFilter(item *PipelineItem) {
	atomic.AddInt64(&p.metrics.FilterProcessed, 1)

	logEntry := item.ParsedLog.SyslogLog

	// --- Filter policy matching ---
	if item.Device != nil {
		result, err := p.filterService.FilterByDevice(item.Device.ID, item.ParsedData)
		if err != nil {
			log.Printf("[pipeline] Filter error for %s: %v", logEntry.LogID, err)
		} else {
			item.FilterResult = result
			logEntry.FilterStatus = result.Action

			if result.Policy != nil {
				logEntry.MatchedFilterPolicyID = &result.Policy.ID
			}

			if item.Trace != nil {
				_ = p.traceService.UpdateTrace(logEntry.LogID, map[string]interface{}{
					"filter_status":     result.Action,
					"matched_policy_id": logEntry.MatchedFilterPolicyID,
					"matched_policy_name": func() string {
						if result.Policy != nil {
							return result.Policy.Name
						}
						return ""
					}(),
				})
				item.Trace.FilterStatus = result.Action
				if result.Policy != nil {
					item.Trace.MatchedPolicyID = &result.Policy.ID
					item.Trace.MatchedPolicyName = result.Policy.Name
				}
			}

			// If action is "drop", stop processing
			if result.Action == "drop" {
				logEntry.FilterStatus = "dropped"
				// Still write to DB for record
				sendOrDrop(p.dbCh, item, &p.metrics.DBDropped, "db")
				return
			}
		}
	} else {
		logEntry.FilterStatus = "no_device"
	}

	// --- Dedup check ---
	isDup, dupCount := p.dedupService.IsDuplicate(logEntry.RawMessage, 60)
	item.Deduped = isDup
	if item.Trace != nil {
		if isDup {
			item.Trace.DedupResult = fmt.Sprintf("duplicate (count=%d)", dupCount)
		} else {
			item.Trace.DedupResult = "unique"
		}
	}

	// --- Aggregation ---
	if logEntry.SourceIP != "" && logEntry.Severity != "" {
		alert, err := p.aggregateService.Aggregate(
			logEntry.SourceIP,
			logEntry.DestinationIP,
			item.extractEventType(),
			logEntry.Severity,
			logEntry.DeviceID,
		)
		if err == nil && alert != nil {
			if item.Trace != nil {
				b, _ := json.Marshal(alert)
				item.Trace.AggregationResult = string(b)
			}
			logEntry.AggregatedAlertID = &alert.ID
		}
	}

	// --- High-frequency IP detection ---
	highFreq, _ := p.highFreqService.Detect(logEntry.SourceIP)
	item.HighFreq = highFreq
	if highFreq {
		_ = p.highFreqService.PersistHighFreqIP(logEntry.SourceIP)
	}

	// --- Desensitization ---
	if item.ParseOK && len(item.ParsedData) > 0 {
		desensitized, err := p.desensitizeService.Desensitize(item.ParsedData)
		if err == nil {
			b, _ := json.Marshal(desensitized)
			logEntry.ParsedData = string(b)
		}
	}

	// Send to DB writer
	sendOrDrop(p.dbCh, item, &p.metrics.DBDropped, "db")

	// Send to push workers (if not deduped)
	if !isDup {
		sendOrDrop(p.pushCh, item, &p.metrics.PushDropped, "push")
	}
}

// =============================================================================
// Stage 3: DB writer (batch insert)
// =============================================================================

// batch buffer for DB writes
var (
	dbBatchMu    sync.Mutex
	dbBatch      []*models.SyslogLog
	dbTraceBatch []*models.LogTraceInfo
)

const dbBatchSize = 100

func (p *Pipeline) dbWriter() {
	defer p.wg.Done()
	defer p.recoverPanic("dbWriter", -1)

	log.Println("[pipeline] DB writer started")

	for {
		select {
		case <-p.stopCh:
			p.drainDB()
			log.Println("[pipeline] DB writer stopped")
			return
		case item := <-p.dbCh:
			p.writeToDB(item)
		}
	}
}

func (p *Pipeline) drainDB() {
	for {
		select {
		case item := <-p.dbCh:
			p.writeToDB(item)
		default:
			p.flushBatch()
			return
		}
	}
}

func (p *Pipeline) writeToDB(item *PipelineItem) {
	dbBatchMu.Lock()
	dbBatch = append(dbBatch, item.ParsedLog.SyslogLog)
	if item.Trace != nil {
		dbTraceBatch = append(dbTraceBatch, item.Trace)
	}
	count := len(dbBatch)
	dbBatchMu.Unlock()

	if count >= dbBatchSize {
		p.flushBatch()
	}
}

func (p *Pipeline) dbFlushTicker() {
	defer p.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.flushBatch()
		}
	}
}

func (p *Pipeline) flushBatch() {
	dbBatchMu.Lock()
	syslogs := dbBatch
	traces := dbTraceBatch
	dbBatch = nil
	dbTraceBatch = nil
	dbBatchMu.Unlock()

	if len(syslogs) == 0 {
		return
	}

	db := database.GetDB()
	if db == nil {
		atomic.AddInt64(&p.metrics.DBErrors, int64(len(syslogs)))
		log.Println("[pipeline] DB not available, dropping batch")
		return
	}

	// Batch insert logs
	if err := db.CreateInBatches(syslogs, dbBatchSize).Error; err != nil {
		atomic.AddInt64(&p.metrics.DBErrors, int64(len(syslogs)))
		log.Printf("[pipeline] DB batch write error: %v", err)
		return
	}

	atomic.AddInt64(&p.metrics.DBWritten, int64(len(syslogs)))

	// Update traces
	for _, trace := range traces {
		if trace != nil {
			if err := db.Save(trace).Error; err != nil {
				atomic.AddInt64(&p.metrics.DBErrors, 1)
			}
		}
	}
}

// =============================================================================
// Stage 4: Push workers
// =============================================================================

func (p *Pipeline) pushWorker(id int) {
	defer p.wg.Done()
	defer p.recoverPanic("pushWorker", id)

	log.Printf("[pipeline] Push worker %d started", id)

	for {
		select {
		case <-p.stopCh:
			p.drainPush(id)
			log.Printf("[pipeline] Push worker %d stopped", id)
			return
		case item := <-p.pushCh:
			p.processPush(item)
		}
	}
}

func (p *Pipeline) drainPush(id int) {
	for {
		select {
		case item := <-p.pushCh:
			p.processPush(item)
		default:
			return
		}
	}
}

func (p *Pipeline) processPush(item *PipelineItem) {
	atomic.AddInt64(&p.metrics.PushProcessed, 1)

	logEntry := item.ParsedLog.SyslogLog

	// Trigger alert processing (will find matching alert rules and push)
	record, err := p.alertService.ProcessAlert(
		logEntry.LogID,
		logEntry.SourceIP,
		logEntry.RawMessage,
		logEntry.DeviceID,
	)
	if err != nil {
		atomic.AddInt64(&p.metrics.PushErrors, 1)
		log.Printf("[pipeline] Alert processing error for %s: %v", logEntry.LogID, err)
	}

	if record != nil && item.Trace != nil {
		b, _ := json.Marshal(record)
		item.Trace.PushResults = string(b)
		_ = p.traceService.UpdateTrace(logEntry.LogID, map[string]interface{}{
			"push_results": item.Trace.PushResults,
		})
	}
}

// =============================================================================
// Helpers
// =============================================================================

// sendOrDrop attempts to send an item to a channel. If the channel is full it
// blocks for up to 100ms then drops the item.
func sendOrDrop[T any](ch chan<- T, item T, counter *int64, stage string) {
	select {
	case ch <- item:
		return
	case <-time.After(100 * time.Millisecond):
		atomic.AddInt64(counter, 1)
		log.Printf("[pipeline] %s stage queue full, dropping item", stage)
		return
	}
}

// recoverPanic catches panics in worker goroutines and logs them.
func (p *Pipeline) recoverPanic(worker string, id int) {
	if r := recover(); r != nil {
		log.Printf("[pipeline] PANIC in %s %d: %v\n%s", worker, id, r, string(debug.Stack()))
	}
}

// extractEventType extracts an event type from parsed data.
func (item *PipelineItem) extractEventType() string {
	if item.ParsedData == nil {
		return ""
	}
	if et, ok := item.ParsedData["event_type"].(string); ok {
		return et
	}
	return ""
}
