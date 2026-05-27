import{a5 as X,be as G,t as r,w,s as y,v as m,a3 as $,by as F,al as p,W as J,a$ as U,aq as Y,N as B,X as S,aV as C,V as N,P as a,aT as Z,a0 as h,bC as g,bo as u,b4 as _,_ as b,Q as T,bj as x,R as P,F as ee,b5 as te,bw as ie,bv as ne}from"./index-qenDAa9H.js";import{g as oe}from"./logs-CDpuFe7r.js";import{_ as re}from"./StatusBadge.vue_vue_type_script_setup_true_lang-DQlW_tiC.js";import{P as le}from"./PageHeader-CWyS-U03.js";import{B as se}from"./Button-Bkv4AwgB.js";import{q as ae,F as O,C as I,A as R,D as ce,L as me,l as de}from"./Scrollbar-DUkKw_1i.js";import{N as ue}from"./Spin-DfeekRdc.js";import{N as pe}from"./Icon-D-79eCSJ.js";import{N as ge}from"./Space-D6xIht3K.js";import{N as V}from"./Card-Dq-y6vVg.js";import"./Tag-DOGr91go.js";import"./_plugin-vue_export-helper-vfH7Hgr3.js";import"./get-slot-DjYhNZAV.js";let j=!1;function ve(){if(ae&&window.CSS&&!j&&(j=!0,"registerProperty"in(window==null?void 0:window.CSS)))try{CSS.registerProperty({name:"--n-color-start",syntax:"<color>",inherits:!1,initialValue:"#0000"}),CSS.registerProperty({name:"--n-color-end",syntax:"<color>",inherits:!1,initialValue:"#0000"})}catch{}}function fe(t){const{textColor3:o,infoColor:l,errorColor:s,successColor:e,warningColor:c,textColor1:v,textColor2:f,railColor:z,fontWeightStrong:n,fontSize:i}=t;return Object.assign(Object.assign({},G),{contentFontSize:i,titleFontWeight:n,circleBorder:`2px solid ${o}`,circleBorderInfo:`2px solid ${l}`,circleBorderError:`2px solid ${s}`,circleBorderSuccess:`2px solid ${e}`,circleBorderWarning:`2px solid ${c}`,iconColor:o,iconColorInfo:l,iconColorError:s,iconColorSuccess:e,iconColorWarning:c,titleTextColor:v,contentTextColor:f,metaTextColor:o,lineColor:z})}const he={common:X,self:fe},A=1.25,xe=r("timeline",`
 position: relative;
 width: 100%;
 display: flex;
 flex-direction: column;
 line-height: ${A};
`,[w("horizontal",`
 flex-direction: row;
 `,[y(">",[r("timeline-item",`
 flex-shrink: 0;
 padding-right: 40px;
 `,[w("dashed-line-type",[y(">",[r("timeline-item-timeline",[m("line",`
 background-image: linear-gradient(90deg, var(--n-color-start), var(--n-color-start) 50%, transparent 50%, transparent 100%);
 background-size: 10px 1px;
 `)])])]),y(">",[r("timeline-item-content",`
 margin-top: calc(var(--n-icon-size) + 12px);
 `,[y(">",[m("meta",`
 margin-top: 6px;
 margin-bottom: unset;
 `)])]),r("timeline-item-timeline",`
 width: 100%;
 height: calc(var(--n-icon-size) + 12px);
 `,[m("line",`
 left: var(--n-icon-size);
 top: calc(var(--n-icon-size) / 2 - 1px);
 right: 0px;
 width: unset;
 height: 2px;
 `)])])])])]),w("right-placement",[r("timeline-item",[r("timeline-item-content",`
 text-align: right;
 margin-right: calc(var(--n-icon-size) + 12px);
 `),r("timeline-item-timeline",`
 width: var(--n-icon-size);
 right: 0;
 `)])]),w("left-placement",[r("timeline-item",[r("timeline-item-content",`
 margin-left: calc(var(--n-icon-size) + 12px);
 `),r("timeline-item-timeline",`
 left: 0;
 `)])]),r("timeline-item",`
 position: relative;
 `,[y("&:last-child",[r("timeline-item-timeline",[m("line",`
 display: none;
 `)]),r("timeline-item-content",[m("meta",`
 margin-bottom: 0;
 `)])]),r("timeline-item-content",[m("title",`
 margin: var(--n-title-margin);
 font-size: var(--n-title-font-size);
 transition: color .3s var(--n-bezier);
 font-weight: var(--n-title-font-weight);
 color: var(--n-title-text-color);
 `),m("content",`
 transition: color .3s var(--n-bezier);
 font-size: var(--n-content-font-size);
 color: var(--n-content-text-color);
 `),m("meta",`
 transition: color .3s var(--n-bezier);
 font-size: 12px;
 margin-top: 6px;
 margin-bottom: 20px;
 color: var(--n-meta-text-color);
 `)]),w("dashed-line-type",[r("timeline-item-timeline",[m("line",`
 --n-color-start: var(--n-line-color);
 transition: --n-color-start .3s var(--n-bezier);
 background-color: transparent;
 background-image: linear-gradient(180deg, var(--n-color-start), var(--n-color-start) 50%, transparent 50%, transparent 100%);
 background-size: 1px 10px;
 `)])]),r("timeline-item-timeline",`
 width: calc(var(--n-icon-size) + 12px);
 position: absolute;
 top: calc(var(--n-title-font-size) * ${A} / 2 - var(--n-icon-size) / 2);
 height: 100%;
 `,[m("circle",`
 border: var(--n-circle-border);
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 width: var(--n-icon-size);
 height: var(--n-icon-size);
 border-radius: var(--n-icon-size);
 box-sizing: border-box;
 `),m("icon",`
 color: var(--n-icon-color);
 font-size: var(--n-icon-size);
 height: var(--n-icon-size);
 width: var(--n-icon-size);
 display: flex;
 align-items: center;
 justify-content: center;
 `),m("line",`
 transition: background-color .3s var(--n-bezier);
 position: absolute;
 top: var(--n-icon-size);
 left: calc(var(--n-icon-size) / 2 - 1px);
 bottom: 0px;
 width: 2px;
 background-color: var(--n-line-color);
 `)])])]),ze=Object.assign(Object.assign({},F.props),{horizontal:Boolean,itemPlacement:{type:String,default:"left"},size:{type:String,default:"medium"},iconSize:Number}),L=J("n-timeline"),be=$({name:"Timeline",props:ze,setup(t,{slots:o}){const{mergedClsPrefixRef:l}=O(t),s=F("Timeline","-timeline",xe,he,t,l);return U(L,{props:t,mergedThemeRef:s,mergedClsPrefixRef:l}),()=>{const{value:e}=l;return p("div",{class:[`${e}-timeline`,t.horizontal&&`${e}-timeline--horizontal`,`${e}-timeline--${t.size}-size`,!t.horizontal&&`${e}-timeline--${t.itemPlacement}-placement`]},o)}}}),Ce={time:[String,Number],title:String,content:String,color:String,lineType:{type:String,default:"default"},type:{type:String,default:"default"}},we=$({name:"TimelineItem",props:Ce,slots:Object,setup(t){const o=Y(L);o||ce("timeline-item","`n-timeline-item` must be placed inside `n-timeline`."),ve();const{inlineThemeDisabled:l}=O(),s=B(()=>{const{props:{size:c,iconSize:v},mergedThemeRef:f}=o,{type:z}=t,{self:{titleTextColor:n,contentTextColor:i,metaTextColor:d,lineColor:k,titleFontWeight:M,contentFontSize:W,[S("iconSize",c)]:D,[S("titleMargin",c)]:E,[S("titleFontSize",c)]:H,[S("circleBorder",z)]:K,[S("iconColor",z)]:q},common:{cubicBezierEaseInOut:Q}}=f.value;return{"--n-bezier":Q,"--n-circle-border":K,"--n-icon-color":q,"--n-content-font-size":W,"--n-content-text-color":i,"--n-line-color":k,"--n-meta-text-color":d,"--n-title-font-size":H,"--n-title-font-weight":M,"--n-title-margin":E,"--n-title-text-color":n,"--n-icon-size":de(v)||D}}),e=l?me("timeline-item",B(()=>{const{props:{size:c,iconSize:v}}=o,{type:f}=t;return`${c[0]}${v||"a"}${f[0]}`}),s,o.props):void 0;return{mergedClsPrefix:o.mergedClsPrefixRef,cssVars:l?void 0:s,themeClass:e==null?void 0:e.themeClass,onRender:e==null?void 0:e.onRender}},render(){const{mergedClsPrefix:t,color:o,onRender:l,$slots:s}=this;return l==null||l(),p("div",{class:[`${t}-timeline-item`,this.themeClass,`${t}-timeline-item--${this.type}-type`,`${t}-timeline-item--${this.lineType}-line-type`],style:this.cssVars},p("div",{class:`${t}-timeline-item-timeline`},p("div",{class:`${t}-timeline-item-timeline__line`}),I(s.icon,e=>e?p("div",{class:`${t}-timeline-item-timeline__icon`,style:{color:o}},e):p("div",{class:`${t}-timeline-item-timeline__circle`,style:{borderColor:o}}))),p("div",{class:`${t}-timeline-item-content`},I(s.header,e=>e||this.title?p("div",{class:`${t}-timeline-item-content__title`},e||this.title):null),p("div",{class:`${t}-timeline-item-content__content`},R(s.default,()=>[this.content])),p("div",{class:`${t}-timeline-item-content__meta`},R(s.footer,()=>[this.time]))))}}),ye={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Se=a("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"48",d:"M244 400L100 256l144-144"},null,-1),$e=a("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"48",d:"M120 256h292"},null,-1),ke=[Se,$e],_e=$({name:"ArrowBackOutline",render:function(o,l){return C(),N("svg",ye,ke)}}),Te={class:"page-container"},Be={style:{"margin-top":"4px",padding:"8px",background:"var(--bg-color)","border-radius":"4px","font-size":"12px","overflow-x":"auto"}},Ee=$({__name:"LogTraceView",setup(t){const o=ne(),l=ie(),s=_(!1),e=_(null),c=_([]),v=B(()=>o.params.id);async function f(){var n,i;if(v.value){s.value=!0;try{const d=await oe(v.value);e.value=((n=d.data)==null?void 0:n.log)||null,c.value=((i=d.data)==null?void 0:i.trace)||z(e.value)}catch{e.value=null,c.value=[]}finally{s.value=!1}}}function z(n){return n?[{stage:"接收",status:"success",detail:`${n.deviceName} (${n.deviceHost})`,time:n.receivedAt},{stage:"解析",status:"success",detail:n.parsedData?"解析完成":"原始消息",time:n.receivedAt},{stage:"过滤",status:"success",detail:"通过过滤策略",time:n.receivedAt},{stage:"去重",status:"success",detail:"无重复",time:n.receivedAt},{stage:"聚合",status:"success",detail:"聚合处理",time:n.receivedAt},{stage:"推送",status:n.pushStatus==="success"?"success":n.pushStatus==="failed"?"error":"warning",detail:n.pushStatus,time:n.receivedAt}]:[]}return Z(()=>{f()}),(n,i)=>(C(),N("div",Te,[h(le,{title:"日志追踪"},{default:g(()=>[h(u(se),{text:"",onClick:i[0]||(i[0]=d=>u(l).push("/logs"))},{icon:g(()=>[h(u(pe),{component:u(_e)},null,8,["component"])]),default:g(()=>[i[1]||(i[1]=b(" 返回 ",-1))]),_:1})]),_:1}),h(u(ue),{show:s.value},{default:g(()=>[e.value?(C(),T(u(V),{key:0,size:"small",style:{"margin-bottom":"16px"}},{default:g(()=>[h(u(ge),{vertical:""},{default:g(()=>[a("div",null,[i[2]||(i[2]=a("strong",null,"日志ID:",-1)),b(" "+x(e.value.id),1)]),a("div",null,[i[3]||(i[3]=a("strong",null,"时间:",-1)),b(" "+x(e.value.receivedAt),1)]),a("div",null,[i[4]||(i[4]=a("strong",null,"设备:",-1)),b(" "+x(e.value.deviceName)+" ("+x(e.value.deviceHost)+")",1)]),a("div",null,[i[5]||(i[5]=a("strong",null,"源IP:",-1)),b(" "+x(e.value.sourceIp),1)]),a("div",null,[i[6]||(i[6]=a("strong",null,"严重程度:",-1)),i[7]||(i[7]=b()),e.value.severity?(C(),T(re,{key:0,status:e.value.severity,type:"severity"},null,8,["status"])):P("",!0)]),a("div",null,[i[8]||(i[8]=a("strong",null,"原始消息:",-1)),a("pre",Be,x(e.value.rawMessage),1)])]),_:1})]),_:1})):P("",!0),h(u(V),{title:"处理链路",size:"small"},{default:g(()=>[h(u(be),null,{default:g(()=>[(C(!0),N(ee,null,te(c.value,(d,k)=>(C(),T(u(we),{key:k,type:d.status==="error"?"error":d.status==="warning"?"warning":"success",title:d.stage,time:d.time},{default:g(()=>[a("div",null,x(d.detail||"-"),1)]),_:2},1032,["type","title","time"]))),128))]),_:1})]),_:1})]),_:1},8,["show"])]))}});export{Ee as default};
