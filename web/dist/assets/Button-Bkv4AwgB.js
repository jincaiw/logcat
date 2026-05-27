import{K as ke,q as K,p as He,C as te,s as Re,c as Fe,b as Ee,G as Ie,F as We,J as Oe,L as je,e as Le,f as re}from"./Scrollbar-DUkKw_1i.js";import{aR as ae,aq as se,N as k,W as le,a$ as Ne,a3 as J,al as x,f as _e,e as De,s as g,G as Ae,t as G,bk as Ge,aM as Ke,b4 as j,L as de,w as C,v as h,x as ne,by as ce,r as Me,X as r,z as D}from"./index-qenDAa9H.js";const ie=le("n-form-item");function Ve(e,{defaultSize:b="medium",mergedSize:v,mergedDisabled:u}={}){const i=se(ie,null);Ne(ie,null);const I=k(v?()=>v(i):()=>{const{size:s}=e;if(s)return s;if(i){const{mergedSize:w}=i;if(w.value!==void 0)return w.value}return b}),H=k(u?()=>u(i):()=>{const{disabled:s}=e;return s!==void 0?s:i?i.disabled.value:!1}),t=k(()=>{const{status:s}=e;return s||(i==null?void 0:i.mergedValidationStatus.value)});return ae(()=>{i&&i.restoreValidation()}),{mergedSizeRef:I,mergedDisabledRef:H,mergedStatusRef:t,nTriggerFormBlur(){i&&i.handleContentBlur()},nTriggerFormChange(){i&&i.handleContentChange()},nTriggerFormFocus(){i&&i.handleContentFocus()},nTriggerFormInput(){i&&i.handleContentInput()}}}const qe=J({name:"FadeInExpandTransition",props:{appear:Boolean,group:Boolean,mode:String,onLeave:Function,onAfterLeave:Function,onAfterEnter:Function,width:Boolean,reverse:Boolean},setup(e,{slots:b}){function v(t){e.width?t.style.maxWidth=`${t.offsetWidth}px`:t.style.maxHeight=`${t.offsetHeight}px`,t.offsetWidth}function u(t){e.width?t.style.maxWidth="0":t.style.maxHeight="0",t.offsetWidth;const{onLeave:s}=e;s&&s()}function i(t){e.width?t.style.maxWidth="":t.style.maxHeight="";const{onAfterLeave:s}=e;s&&s()}function I(t){if(t.style.transition="none",e.width){const s=t.offsetWidth;t.style.maxWidth="0",t.offsetWidth,t.style.transition="",t.style.maxWidth=`${s}px`}else if(e.reverse)t.style.maxHeight=`${t.offsetHeight}px`,t.offsetHeight,t.style.transition="",t.style.maxHeight="0";else{const s=t.offsetHeight;t.style.maxHeight="0",t.offsetWidth,t.style.transition="",t.style.maxHeight=`${s}px`}t.offsetWidth}function H(t){var s;e.width?t.style.maxWidth="":e.reverse||(t.style.maxHeight=""),(s=e.onAfterEnter)===null||s===void 0||s.call(e)}return()=>{const{group:t,width:s,appear:w,mode:W}=e,O=t?_e:De,L={name:s?"fade-in-width-expand-transition":"fade-in-height-expand-transition",appear:w,onEnter:I,onAfterEnter:H,onBeforeLeave:v,onLeave:u,onAfterLeave:i};return t||(L.mode=W),x(O,L,b)}}}),{cubicBezierEaseInOut:z}=Ae;function Qe({duration:e=".2s",delay:b=".1s"}={}){return[g("&.fade-in-width-expand-transition-leave-from, &.fade-in-width-expand-transition-enter-to",{opacity:1}),g("&.fade-in-width-expand-transition-leave-to, &.fade-in-width-expand-transition-enter-from",`
 opacity: 0!important;
 margin-left: 0!important;
 margin-right: 0!important;
 `),g("&.fade-in-width-expand-transition-leave-active",`
 overflow: hidden;
 transition:
 opacity ${e} ${z},
 max-width ${e} ${z} ${b},
 margin-left ${e} ${z} ${b},
 margin-right ${e} ${z} ${b};
 `),g("&.fade-in-width-expand-transition-enter-active",`
 overflow: hidden;
 transition:
 opacity ${e} ${z} ${b},
 max-width ${e} ${z},
 margin-left ${e} ${z},
 margin-right ${e} ${z};
 `)]}const Xe=G("base-wave",`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
`),Ye=J({name:"BaseWave",props:{clsPrefix:{type:String,required:!0}},setup(e){ke("-base-wave",Xe,Ge(e,"clsPrefix"));const b=j(null),v=j(!1);let u=null;return ae(()=>{u!==null&&window.clearTimeout(u)}),{active:v,selfRef:b,play(){u!==null&&(window.clearTimeout(u),v.value=!1,u=null),Ke(()=>{var i;(i=b.value)===null||i===void 0||i.offsetHeight,v.value=!0,u=window.setTimeout(()=>{v.value=!1,u=null},1e3)})}}},render(){const{clsPrefix:e}=this;return x("div",{ref:"selfRef","aria-hidden":!0,class:[`${e}-base-wave`,this.active&&`${e}-base-wave--active`]})}}),Je=K&&"chrome"in window;K&&navigator.userAgent.includes("Firefox");const Ue=K&&navigator.userAgent.includes("Safari")&&!Je;function P(e){return de(e,[255,255,255,.16])}function A(e){return de(e,[0,0,0,.12])}const Ze=le("n-button-group"),eo=g([G("button",`
 margin: 0;
 font-weight: var(--n-font-weight);
 line-height: 1;
 font-family: inherit;
 padding: var(--n-padding);
 height: var(--n-height);
 font-size: var(--n-font-size);
 border-radius: var(--n-border-radius);
 color: var(--n-text-color);
 background-color: var(--n-color);
 width: var(--n-width);
 white-space: nowrap;
 outline: none;
 position: relative;
 z-index: auto;
 border: none;
 display: inline-flex;
 flex-wrap: nowrap;
 flex-shrink: 0;
 align-items: center;
 justify-content: center;
 user-select: none;
 -webkit-user-select: none;
 text-align: center;
 cursor: pointer;
 text-decoration: none;
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[C("color",[h("border",{borderColor:"var(--n-border-color)"}),C("disabled",[h("border",{borderColor:"var(--n-border-color-disabled)"})]),ne("disabled",[g("&:focus",[h("state-border",{borderColor:"var(--n-border-color-focus)"})]),g("&:hover",[h("state-border",{borderColor:"var(--n-border-color-hover)"})]),g("&:active",[h("state-border",{borderColor:"var(--n-border-color-pressed)"})]),C("pressed",[h("state-border",{borderColor:"var(--n-border-color-pressed)"})])])]),C("disabled",{backgroundColor:"var(--n-color-disabled)",color:"var(--n-text-color-disabled)"},[h("border",{border:"var(--n-border-disabled)"})]),ne("disabled",[g("&:focus",{backgroundColor:"var(--n-color-focus)",color:"var(--n-text-color-focus)"},[h("state-border",{border:"var(--n-border-focus)"})]),g("&:hover",{backgroundColor:"var(--n-color-hover)",color:"var(--n-text-color-hover)"},[h("state-border",{border:"var(--n-border-hover)"})]),g("&:active",{backgroundColor:"var(--n-color-pressed)",color:"var(--n-text-color-pressed)"},[h("state-border",{border:"var(--n-border-pressed)"})]),C("pressed",{backgroundColor:"var(--n-color-pressed)",color:"var(--n-text-color-pressed)"},[h("state-border",{border:"var(--n-border-pressed)"})])]),C("loading","cursor: wait;"),G("base-wave",`
 pointer-events: none;
 top: 0;
 right: 0;
 bottom: 0;
 left: 0;
 animation-iteration-count: 1;
 animation-duration: var(--n-ripple-duration);
 animation-timing-function: var(--n-bezier-ease-out), var(--n-bezier-ease-out);
 `,[C("active",{zIndex:1,animationName:"button-wave-spread, button-wave-opacity"})]),K&&"MozBoxSizing"in document.createElement("div").style?g("&::moz-focus-inner",{border:0}):null,h("border, state-border",`
 position: absolute;
 left: 0;
 top: 0;
 right: 0;
 bottom: 0;
 border-radius: inherit;
 transition: border-color .3s var(--n-bezier);
 pointer-events: none;
 `),h("border",`
 border: var(--n-border);
 `),h("state-border",`
 border: var(--n-border);
 border-color: #0000;
 z-index: 1;
 `),h("icon",`
 margin: var(--n-icon-margin);
 margin-left: 0;
 height: var(--n-icon-size);
 width: var(--n-icon-size);
 max-width: var(--n-icon-size);
 font-size: var(--n-icon-size);
 position: relative;
 flex-shrink: 0;
 `,[G("icon-slot",`
 height: var(--n-icon-size);
 width: var(--n-icon-size);
 position: absolute;
 left: 0;
 top: 50%;
 transform: translateY(-50%);
 display: flex;
 align-items: center;
 justify-content: center;
 `,[He({top:"50%",originalTransform:"translateY(-50%)"})]),Qe()]),h("content",`
 display: flex;
 align-items: center;
 flex-wrap: nowrap;
 min-width: 0;
 `,[g("~",[h("icon",{margin:"var(--n-icon-margin)",marginRight:0})])]),C("block",`
 display: flex;
 width: 100%;
 `),C("dashed",[h("border, state-border",{borderStyle:"dashed !important"})]),C("disabled",{cursor:"not-allowed",opacity:"var(--n-opacity-disabled)"})]),g("@keyframes button-wave-spread",{from:{boxShadow:"0 0 0.5px 0 var(--n-ripple-color)"},to:{boxShadow:"0 0 0.5px 4.5px var(--n-ripple-color)"}}),g("@keyframes button-wave-opacity",{from:{opacity:"var(--n-wave-opacity)"},to:{opacity:0}})]),oo=Object.assign(Object.assign({},ce.props),{color:String,textColor:String,text:Boolean,block:Boolean,loading:Boolean,disabled:Boolean,circle:Boolean,size:String,ghost:Boolean,round:Boolean,secondary:Boolean,tertiary:Boolean,quaternary:Boolean,strong:Boolean,focusable:{type:Boolean,default:!0},keyboard:{type:Boolean,default:!0},tag:{type:String,default:"button"},type:{type:String,default:"default"},dashed:Boolean,renderIcon:Function,iconPlacement:{type:String,default:"left"},attrType:{type:String,default:"button"},bordered:{type:Boolean,default:!0},onClick:[Function,Array],nativeFocusBehavior:{type:Boolean,default:!Ue},spinProps:Object}),to=J({name:"Button",props:oo,slots:Object,setup(e){const b=j(null),v=j(null),u=j(!1),i=Ie(()=>!e.quaternary&&!e.tertiary&&!e.secondary&&!e.text&&(!e.color||e.ghost||e.dashed)&&e.bordered),I=se(Ze,{}),{inlineThemeDisabled:H,mergedClsPrefixRef:t,mergedRtlRef:s,mergedComponentPropsRef:w}=We(e),{mergedSizeRef:W}=Ve({},{defaultSize:"medium",mergedSize:n=>{var f,y;const{size:o}=e;if(o)return o;const{size:F}=I;if(F)return F;const{mergedSize:B}=n||{};if(B)return B.value;const E=(y=(f=w==null?void 0:w.value)===null||f===void 0?void 0:f.Button)===null||y===void 0?void 0:y.size;return E||"medium"}}),O=k(()=>e.focusable&&!e.disabled),L=n=>{var f;O.value||n.preventDefault(),!e.nativeFocusBehavior&&(n.preventDefault(),!e.disabled&&O.value&&((f=b.value)===null||f===void 0||f.focus({preventScroll:!0})))},ue=n=>{var f;if(!e.disabled&&!e.loading){const{onClick:y}=e;y&&Le(y,n),e.text||(f=v.value)===null||f===void 0||f.play()}},fe=n=>{switch(n.key){case"Enter":if(!e.keyboard)return;u.value=!1}},he=n=>{switch(n.key){case"Enter":if(!e.keyboard||e.loading){n.preventDefault();return}u.value=!0}},be=()=>{u.value=!1},ve=ce("Button","-button",eo,Me,e,t),me=Oe("Button",s,t),U=k(()=>{const n=ve.value,{common:{cubicBezierEaseInOut:f,cubicBezierEaseOut:y},self:o}=n,{rippleDuration:F,opacityDisabled:B,fontWeight:E,fontWeightStrong:M}=o,p=W.value,{dashed:V,type:S,ghost:q,text:$,color:l,round:Z,circle:Q,textColor:T,secondary:ge,tertiary:ee,quaternary:xe,strong:ye}=e,pe={"--n-font-weight":ye?M:E};let d={"--n-color":"initial","--n-color-hover":"initial","--n-color-pressed":"initial","--n-color-focus":"initial","--n-color-disabled":"initial","--n-ripple-color":"initial","--n-text-color":"initial","--n-text-color-hover":"initial","--n-text-color-pressed":"initial","--n-text-color-focus":"initial","--n-text-color-disabled":"initial"};const N=S==="tertiary",oe=S==="default",a=N?"default":S;if($){const c=T||l;d={"--n-color":"#0000","--n-color-hover":"#0000","--n-color-pressed":"#0000","--n-color-focus":"#0000","--n-color-disabled":"#0000","--n-ripple-color":"#0000","--n-text-color":c||o[r("textColorText",a)],"--n-text-color-hover":c?P(c):o[r("textColorTextHover",a)],"--n-text-color-pressed":c?A(c):o[r("textColorTextPressed",a)],"--n-text-color-focus":c?P(c):o[r("textColorTextHover",a)],"--n-text-color-disabled":c||o[r("textColorTextDisabled",a)]}}else if(q||V){const c=T||l;d={"--n-color":"#0000","--n-color-hover":"#0000","--n-color-pressed":"#0000","--n-color-focus":"#0000","--n-color-disabled":"#0000","--n-ripple-color":l||o[r("rippleColor",a)],"--n-text-color":c||o[r("textColorGhost",a)],"--n-text-color-hover":c?P(c):o[r("textColorGhostHover",a)],"--n-text-color-pressed":c?A(c):o[r("textColorGhostPressed",a)],"--n-text-color-focus":c?P(c):o[r("textColorGhostHover",a)],"--n-text-color-disabled":c||o[r("textColorGhostDisabled",a)]}}else if(ge){const c=oe?o.textColor:N?o.textColorTertiary:o[r("color",a)],m=l||c,_=S!=="default"&&S!=="tertiary";d={"--n-color":_?D(m,{alpha:Number(o.colorOpacitySecondary)}):o.colorSecondary,"--n-color-hover":_?D(m,{alpha:Number(o.colorOpacitySecondaryHover)}):o.colorSecondaryHover,"--n-color-pressed":_?D(m,{alpha:Number(o.colorOpacitySecondaryPressed)}):o.colorSecondaryPressed,"--n-color-focus":_?D(m,{alpha:Number(o.colorOpacitySecondaryHover)}):o.colorSecondaryHover,"--n-color-disabled":o.colorSecondary,"--n-ripple-color":"#0000","--n-text-color":m,"--n-text-color-hover":m,"--n-text-color-pressed":m,"--n-text-color-focus":m,"--n-text-color-disabled":m}}else if(ee||xe){const c=oe?o.textColor:N?o.textColorTertiary:o[r("color",a)],m=l||c;ee?(d["--n-color"]=o.colorTertiary,d["--n-color-hover"]=o.colorTertiaryHover,d["--n-color-pressed"]=o.colorTertiaryPressed,d["--n-color-focus"]=o.colorSecondaryHover,d["--n-color-disabled"]=o.colorTertiary):(d["--n-color"]=o.colorQuaternary,d["--n-color-hover"]=o.colorQuaternaryHover,d["--n-color-pressed"]=o.colorQuaternaryPressed,d["--n-color-focus"]=o.colorQuaternaryHover,d["--n-color-disabled"]=o.colorQuaternary),d["--n-ripple-color"]="#0000",d["--n-text-color"]=m,d["--n-text-color-hover"]=m,d["--n-text-color-pressed"]=m,d["--n-text-color-focus"]=m,d["--n-text-color-disabled"]=m}else d={"--n-color":l||o[r("color",a)],"--n-color-hover":l?P(l):o[r("colorHover",a)],"--n-color-pressed":l?A(l):o[r("colorPressed",a)],"--n-color-focus":l?P(l):o[r("colorFocus",a)],"--n-color-disabled":l||o[r("colorDisabled",a)],"--n-ripple-color":l||o[r("rippleColor",a)],"--n-text-color":T||(l?o.textColorPrimary:N?o.textColorTertiary:o[r("textColor",a)]),"--n-text-color-hover":T||(l?o.textColorHoverPrimary:o[r("textColorHover",a)]),"--n-text-color-pressed":T||(l?o.textColorPressedPrimary:o[r("textColorPressed",a)]),"--n-text-color-focus":T||(l?o.textColorFocusPrimary:o[r("textColorFocus",a)]),"--n-text-color-disabled":T||(l?o.textColorDisabledPrimary:o[r("textColorDisabled",a)])};let X={"--n-border":"initial","--n-border-hover":"initial","--n-border-pressed":"initial","--n-border-focus":"initial","--n-border-disabled":"initial"};$?X={"--n-border":"none","--n-border-hover":"none","--n-border-pressed":"none","--n-border-focus":"none","--n-border-disabled":"none"}:X={"--n-border":o[r("border",a)],"--n-border-hover":o[r("borderHover",a)],"--n-border-pressed":o[r("borderPressed",a)],"--n-border-focus":o[r("borderFocus",a)],"--n-border-disabled":o[r("borderDisabled",a)]};const{[r("height",p)]:Y,[r("fontSize",p)]:Ce,[r("padding",p)]:we,[r("paddingRound",p)]:$e,[r("iconSize",p)]:ze,[r("borderRadius",p)]:Be,[r("iconMargin",p)]:Se,waveOpacity:Te}=o,Pe={"--n-width":Q&&!$?Y:"initial","--n-height":$?"initial":Y,"--n-font-size":Ce,"--n-padding":Q||$?"initial":Z?$e:we,"--n-icon-size":ze,"--n-icon-margin":Se,"--n-border-radius":$?"initial":Q||Z?Y:Be};return Object.assign(Object.assign(Object.assign(Object.assign({"--n-bezier":f,"--n-bezier-ease-out":y,"--n-ripple-duration":F,"--n-opacity-disabled":B,"--n-wave-opacity":Te},pe),d),X),Pe)}),R=H?je("button",k(()=>{let n="";const{dashed:f,type:y,ghost:o,text:F,color:B,round:E,circle:M,textColor:p,secondary:V,tertiary:S,quaternary:q,strong:$}=e;f&&(n+="a"),o&&(n+="b"),F&&(n+="c"),E&&(n+="d"),M&&(n+="e"),V&&(n+="f"),S&&(n+="g"),q&&(n+="h"),$&&(n+="i"),B&&(n+=`j${re(B)}`),p&&(n+=`k${re(p)}`);const{value:l}=W;return n+=`l${l[0]}`,n+=`m${y[0]}`,n}),U,e):void 0;return{selfElRef:b,waveElRef:v,mergedClsPrefix:t,mergedFocusable:O,mergedSize:W,showBorder:i,enterPressed:u,rtlEnabled:me,handleMousedown:L,handleKeydown:he,handleBlur:be,handleKeyup:fe,handleClick:ue,customColorCssVars:k(()=>{const{color:n}=e;if(!n)return null;const f=P(n);return{"--n-border-color":n,"--n-border-color-hover":f,"--n-border-color-pressed":A(n),"--n-border-color-focus":f,"--n-border-color-disabled":n}}),cssVars:H?void 0:U,themeClass:R==null?void 0:R.themeClass,onRender:R==null?void 0:R.onRender}},render(){const{mergedClsPrefix:e,tag:b,onRender:v}=this;v==null||v();const u=te(this.$slots.default,i=>i&&x("span",{class:`${e}-button__content`},i));return x(b,{ref:"selfElRef",class:[this.themeClass,`${e}-button`,`${e}-button--${this.type}-type`,`${e}-button--${this.mergedSize}-type`,this.rtlEnabled&&`${e}-button--rtl`,this.disabled&&`${e}-button--disabled`,this.block&&`${e}-button--block`,this.enterPressed&&`${e}-button--pressed`,!this.text&&this.dashed&&`${e}-button--dashed`,this.color&&`${e}-button--color`,this.secondary&&`${e}-button--secondary`,this.loading&&`${e}-button--loading`,this.ghost&&`${e}-button--ghost`],tabindex:this.mergedFocusable?0:-1,type:this.attrType,style:this.cssVars,disabled:this.disabled,onClick:this.handleClick,onBlur:this.handleBlur,onMousedown:this.handleMousedown,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},this.iconPlacement==="right"&&u,x(qe,{width:!0},{default:()=>te(this.$slots.icon,i=>(this.loading||this.renderIcon||i)&&x("span",{class:`${e}-button__icon`,style:{margin:Re(this.$slots.default)?"0":""}},x(Fe,null,{default:()=>this.loading?x(Ee,Object.assign({clsPrefix:e,key:"loading",class:`${e}-icon-slot`,strokeWidth:20},this.spinProps)):x("div",{key:"icon",class:`${e}-icon-slot`,role:"none"},this.renderIcon?this.renderIcon():i)})))}),this.iconPlacement==="left"&&u,this.text?null:x(Ye,{ref:"waveElRef",clsPrefix:e}),this.showBorder?x("div",{"aria-hidden":!0,class:`${e}-button__border`,style:this.customColorCssVars}):null,this.showBorder?x("div",{"aria-hidden":!0,class:`${e}-button__state-border`,style:this.customColorCssVars}):null)}}),io=to;export{to as B,qe as N,io as X,ie as f,Ue as i,Ve as u};
