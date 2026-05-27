import{F as j,J as F,L as H}from"./Scrollbar-DUkKw_1i.js";import{s as r,t as i,w as n,x as D,as as I,at as J,a3 as K,al as W,by as h,N as s,bh as X,X as b}from"./index-qenDAa9H.js";const q=r([i("table",`
 font-size: var(--n-font-size);
 font-variant-numeric: tabular-nums;
 line-height: var(--n-line-height);
 width: 100%;
 border-radius: var(--n-border-radius) var(--n-border-radius) 0 0;
 text-align: left;
 border-collapse: separate;
 border-spacing: 0;
 overflow: hidden;
 background-color: var(--n-td-color);
 border-color: var(--n-merged-border-color);
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 --n-merged-border-color: var(--n-border-color);
 `,[r("th",`
 white-space: nowrap;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 text-align: inherit;
 padding: var(--n-th-padding);
 vertical-align: inherit;
 text-transform: none;
 border: 0px solid var(--n-merged-border-color);
 font-weight: var(--n-th-font-weight);
 color: var(--n-th-text-color);
 background-color: var(--n-th-color);
 border-bottom: 1px solid var(--n-merged-border-color);
 border-right: 1px solid var(--n-merged-border-color);
 `,[r("&:last-child",`
 border-right: 0px solid var(--n-merged-border-color);
 `)]),r("td",`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 padding: var(--n-td-padding);
 color: var(--n-td-text-color);
 background-color: var(--n-td-color);
 border: 0px solid var(--n-merged-border-color);
 border-right: 1px solid var(--n-merged-border-color);
 border-bottom: 1px solid var(--n-merged-border-color);
 `,[r("&:last-child",`
 border-right: 0px solid var(--n-merged-border-color);
 `)]),n("bordered",`
 border: 1px solid var(--n-merged-border-color);
 border-radius: var(--n-border-radius);
 `,[r("tr",[r("&:last-child",[r("td",`
 border-bottom: 0 solid var(--n-merged-border-color);
 `)])])]),n("single-line",[r("th",`
 border-right: 0px solid var(--n-merged-border-color);
 `),r("td",`
 border-right: 0px solid var(--n-merged-border-color);
 `)]),n("single-column",[r("tr",[r("&:not(:last-child)",[r("td",`
 border-bottom: 0px solid var(--n-merged-border-color);
 `)])])]),n("striped",[r("tr:nth-of-type(even)",[r("td","background-color: var(--n-td-color-striped)")])]),D("bottom-bordered",[r("tr",[r("&:last-child",[r("td",`
 border-bottom: 0px solid var(--n-merged-border-color);
 `)])])])]),I(i("table",`
 background-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `,[r("th",`
 background-color: var(--n-th-color-modal);
 `),r("td",`
 background-color: var(--n-td-color-modal);
 `)])),J(i("table",`
 background-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `,[r("th",`
 background-color: var(--n-th-color-popover);
 `),r("td",`
 background-color: var(--n-td-color-popover);
 `)]))]),A=Object.assign(Object.assign({},h.props),{bordered:{type:Boolean,default:!0},bottomBordered:{type:Boolean,default:!0},singleLine:{type:Boolean,default:!0},striped:Boolean,singleColumn:Boolean,size:String}),U=K({name:"Table",props:A,setup(e){const{mergedClsPrefixRef:o,inlineThemeDisabled:c,mergedRtlRef:m,mergedComponentPropsRef:a}=j(e),v=s(()=>{var d,l;return e.size||((l=(d=a==null?void 0:a.value)===null||d===void 0?void 0:d.Table)===null||l===void 0?void 0:l.size)||"medium"}),p=h("Table","-table",q,X,e,o),u=F("Table",m,o),g=s(()=>{const d=v.value,{self:{borderColor:l,tdColor:f,tdColorModal:x,tdColorPopover:C,thColor:z,thColorModal:P,thColorPopover:k,thTextColor:R,tdTextColor:T,borderRadius:B,thFontWeight:w,lineHeight:y,borderColorModal:$,borderColorPopover:M,tdColorStriped:S,tdColorStripedModal:E,tdColorStripedPopover:L,[b("fontSize",d)]:N,[b("tdPadding",d)]:O,[b("thPadding",d)]:V},common:{cubicBezierEaseInOut:_}}=p.value;return{"--n-bezier":_,"--n-td-color":f,"--n-td-color-modal":x,"--n-td-color-popover":C,"--n-td-text-color":T,"--n-border-color":l,"--n-border-color-modal":$,"--n-border-color-popover":M,"--n-border-radius":B,"--n-font-size":N,"--n-th-color":z,"--n-th-color-modal":P,"--n-th-color-popover":k,"--n-th-font-weight":w,"--n-th-text-color":R,"--n-line-height":y,"--n-td-padding":O,"--n-th-padding":V,"--n-td-color-striped":S,"--n-td-color-striped-modal":E,"--n-td-color-striped-popover":L}}),t=c?H("table",s(()=>v.value[0]),g,e):void 0;return{rtlEnabled:u,mergedClsPrefix:o,cssVars:c?void 0:g,themeClass:t==null?void 0:t.themeClass,onRender:t==null?void 0:t.onRender}},render(){var e;const{mergedClsPrefix:o}=this;return(e=this.onRender)===null||e===void 0||e.call(this),W("table",{class:[`${o}-table`,this.themeClass,{[`${o}-table--rtl`]:this.rtlEnabled,[`${o}-table--bottom-bordered`]:this.bottomBordered,[`${o}-table--bordered`]:this.bordered,[`${o}-table--single-line`]:this.singleLine,[`${o}-table--single-column`]:this.singleColumn,[`${o}-table--striped`]:this.striped}],style:this.cssVars},this.$slots)}});export{U as N};
