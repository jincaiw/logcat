# User Guide

## 1. Login

Open the web UI and log in.

Default account:

```text
admin / admin123
```

Change the password in **Profile** after first login.

## 2. Configure devices

Go to **Devices** and add the Syslog source devices. Devices help log attribution and policy matching.

Recommended fields:

- Name
- IP address
- Device type
- Parse template
- Group

## 3. Configure parse templates

Go to **Parse Templates**.

Supported parse types:

- JSON
- Syslog + JSON
- Delimiter
- Key-value delimiter
- Regex
- Key-value

Use the built-in preview tool to verify parsed fields before enabling policies.

## 4. Configure filter policies

Go to **Filter Policies** and define matching rules.

A policy can:

- keep or discard logs,
- match multiple conditions,
- use AND/OR logic,
- enable deduplication,
- trigger alert rules.

## 5. Configure notification channels

Go to **Notifications**.

Supported channels:

- Feishu
- Email
- Syslog forwarding

Use **Test Send** to verify connectivity before using a channel in production.

## 6. Create output templates

In **Notifications > Output Templates**, create reusable message formats for alert notifications.

Templates support variable replacement, for example:

```text
Alert: {{alertName}}
Source: {{sourceIp}}
Severity: {{severity}}
```

## 7. Start Syslog service

Use the sidebar service control or **Settings** to start the Syslog receiver.

Default Syslog port: `5140`.

Configure your devices to send Syslog to:

```text
<logcat-server-ip>:5140
```

## 8. Search logs

Go to **Logs** to search received logs by keyword, device, and time range.

## 9. Analyze statistics

Go to **Statistics** for field distribution and Top-N analysis.

## 10. Production checklist

Before production use:

- Change the default admin password.
- Mount or back up the data directory.
- Verify firewall rules for ports `8080` and `5140`.
- Test notification channels.
- Test at least one real device log sample.
- Enable log retention policies according to storage capacity.
