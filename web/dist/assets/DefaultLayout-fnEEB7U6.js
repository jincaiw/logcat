import{a3 as S,al as d,W as ee,t as v,as as uo,s as k,at as ho,v as m,aq as Y,by as q,j as vo,N as x,b4 as F,aT as je,bB as ge,aR as mo,bA as po,X as fo,w as N,m as go,a$ as Q,bk as ae,aU as bo,$ as xo,bb as Co,a5 as wo,L as Oe,x as le,F as Ae,aJ as ko,aH as yo,aV as j,V as E,P as z,Z as He,bu as _o,br as ke,Q as J,bo as B,bw as ye,bv as Le,a0 as K,bC as D,b5 as zo,_ as So,bj as Io,bs as Ro,b8 as $o,b6 as Po,b7 as Oo}from"./index-qenDAa9H.js";import{f as To}from"./fade-in-height-expand.cssr-Dro2CEvB.js";import{r as oe,V as Mo,a as ve}from"./create-CSfaedQP.js";import{A as Ee,C as No,V as Fe,F as te,L as re,f as Bo,q as jo,S as Ke,I as Ve,a as De,l as me,e as W,u as _e,G as be}from"./Scrollbar-DUkKw_1i.js";import{N as Ao}from"./Tooltip-D2F0EjKD.js";import{k as pe}from"./Popover-BJ4Ic_2j.js";import{C as Ho,N as Ue}from"./Dropdown-DJslDNZ9.js";import{N as Lo,B as fe}from"./Button-Bkv4AwgB.js";import{u as xe,a as Eo}from"./get-K_rslPNJ.js";import{u as Fo,_ as he}from"./_plugin-vue_export-helper-vfH7Hgr3.js";import{W as Ko,D as Vo}from"./WarningOutline-uUdoClS3.js";import{N as Z}from"./Icon-D-79eCSJ.js";import{N as Do}from"./Space-D6xIht3K.js";import{i as Uo,o as Go}from"./utils-BbZx3gk3.js";import{t as qo}from"./Tag-DOGr91go.js";import"./utils-DBY6nFc7.js";import"./get-slot-DjYhNZAV.js";const Wo=S({name:"ChevronDownFilled",render(){return d("svg",{viewBox:"0 0 16 16",fill:"none",xmlns:"http://www.w3.org/2000/svg"},d("path",{d:"M3.20041 5.73966C3.48226 5.43613 3.95681 5.41856 4.26034 5.70041L8 9.22652L11.7397 5.70041C12.0432 5.41856 12.5177 5.43613 12.7996 5.73966C13.0815 6.0432 13.0639 6.51775 12.7603 6.7996L8.51034 10.7996C8.22258 11.0668 7.77743 11.0668 7.48967 10.7996L3.23966 6.7996C2.93613 6.51775 2.91856 6.0432 3.20041 5.73966Z",fill:"currentColor"}))}}),Yo=ee("n-avatar-group"),Xo=v("avatar",`
 width: var(--n-merged-size);
 height: var(--n-merged-size);
 color: #FFF;
 font-size: var(--n-font-size);
 display: inline-flex;
 position: relative;
 overflow: hidden;
 text-align: center;
 border: var(--n-border);
 border-radius: var(--n-border-radius);
 --n-merged-color: var(--n-color);
 background-color: var(--n-merged-color);
 transition:
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
`,[uo(k("&","--n-merged-color: var(--n-color-modal);")),ho(k("&","--n-merged-color: var(--n-color-popover);")),k("img",`
 width: 100%;
 height: 100%;
 `),m("text",`
 white-space: nowrap;
 display: inline-block;
 position: absolute;
 left: 50%;
 top: 50%;
 `),v("icon",`
 vertical-align: bottom;
 font-size: calc(var(--n-merged-size) - 6px);
 `),m("text","line-height: 1.25")]),Zo=Object.assign(Object.assign({},q.props),{size:[String,Number],src:String,circle:{type:Boolean,default:void 0},objectFit:String,round:{type:Boolean,default:void 0},bordered:{type:Boolean,default:void 0},onError:Function,fallbackSrc:String,intersectionObserverOptions:Object,lazy:Boolean,onLoad:Function,renderPlaceholder:Function,renderFallback:Function,imgProps:Object,color:String}),Jo=S({name:"Avatar",props:Zo,slots:Object,setup(e){const{mergedClsPrefixRef:o,inlineThemeDisabled:t}=te(e),i=F(!1);let s=null;const n=F(null),u=F(null),c=()=>{const{value:g}=n;if(g&&(s===null||s!==g.innerHTML)){s=g.innerHTML;const{value:O}=u;if(O){const{offsetWidth:H,offsetHeight:T}=O,{offsetWidth:w,offsetHeight:R}=g,U=.9,G=Math.min(H/w*U,T/R*U,1);g.style.transform=`translateX(-50%) translateY(-50%) scale(${G})`}}},l=Y(Yo,null),h=x(()=>{const{size:g}=e;if(g)return g;const{size:O}=l||{};return O||"medium"}),M=q("Avatar","-avatar",Xo,vo,e,o),P=Y(qo,null),p=x(()=>{if(l)return!0;const{round:g,circle:O}=e;return g!==void 0||O!==void 0?g||O:P?P.roundRef.value:!1}),y=x(()=>l?!0:e.bordered||!1),C=x(()=>{const g=h.value,O=p.value,H=y.value,{color:T}=e,{self:{borderRadius:w,fontSize:R,color:U,border:G,colorModal:ne,colorPopover:V},common:{cubicBezierEaseInOut:ce}}=M.value;let ie;return typeof g=="number"?ie=`${g}px`:ie=M.value.self[fo("height",g)],{"--n-font-size":R,"--n-border":H?G:"none","--n-border-radius":O?"50%":w,"--n-color":T||U,"--n-color-modal":T||ne,"--n-color-popover":T||V,"--n-bezier":ce,"--n-merged-size":`var(--n-avatar-size-override, ${ie})`}}),f=t?re("avatar",x(()=>{const g=h.value,O=p.value,H=y.value,{color:T}=e;let w="";return g&&(typeof g=="number"?w+=`a${g}`:w+=g[0]),O&&(w+="b"),H&&(w+="c"),T&&(w+=Bo(T)),w}),C,e):void 0,I=F(!e.lazy);je(()=>{if(e.lazy&&e.intersectionObserverOptions){let g;const O=ge(()=>{g==null||g(),g=void 0,e.lazy&&(g=Go(u.value,e.intersectionObserverOptions,I))});mo(()=>{O(),g==null||g()})}}),po(()=>{var g;return e.src||((g=e.imgProps)===null||g===void 0?void 0:g.src)},()=>{i.value=!1});const A=F(!e.lazy);return{textRef:n,selfRef:u,mergedRoundRef:p,mergedClsPrefix:o,fitTextTransform:c,cssVars:t?void 0:C,themeClass:f==null?void 0:f.themeClass,onRender:f==null?void 0:f.onRender,hasLoadError:i,shouldStartLoading:I,loaded:A,mergedOnError:g=>{if(!I.value)return;i.value=!0;const{onError:O,imgProps:{onError:H}={}}=e;O==null||O(g),H==null||H(g)},mergedOnLoad:g=>{const{onLoad:O,imgProps:{onLoad:H}={}}=e;O==null||O(g),H==null||H(g),A.value=!0}}},render(){var e,o;const{$slots:t,src:i,mergedClsPrefix:s,lazy:n,onRender:u,loaded:c,hasLoadError:l,imgProps:h={}}=this;u==null||u();let M;const P=!c&&!l&&(this.renderPlaceholder?this.renderPlaceholder():(o=(e=this.$slots).placeholder)===null||o===void 0?void 0:o.call(e));return this.hasLoadError?M=this.renderFallback?this.renderFallback():Ee(t.fallback,()=>[d("img",{src:this.fallbackSrc,style:{objectFit:this.objectFit}})]):M=No(t.default,p=>{if(p)return d(Fe,{onResize:this.fitTextTransform},{default:()=>d("span",{ref:"textRef",class:`${s}-avatar__text`},p)});if(i||h.src){const y=this.src||h.src;return d("img",Object.assign(Object.assign({},h),{loading:Uo&&!this.intersectionObserverOptions&&n?"lazy":"eager",src:n&&this.intersectionObserverOptions?this.shouldStartLoading?y:void 0:y,"data-image-src":y,onLoad:this.mergedOnLoad,onError:this.mergedOnError,style:[h.style||"",{objectFit:this.objectFit},P?{height:"0",width:"0",visibility:"hidden",position:"absolute"}:""]}))}}),d("span",{ref:"selfRef",class:[`${s}-avatar`,this.themeClass],style:this.cssVars},M,n&&P)}}),Qo=v("breadcrumb",`
 white-space: nowrap;
 cursor: default;
 line-height: var(--n-item-line-height);
`,[k("ul",`
 list-style: none;
 padding: 0;
 margin: 0;
 `),k("a",`
 color: inherit;
 text-decoration: inherit;
 `),v("breadcrumb-item",`
 font-size: var(--n-font-size);
 transition: color .3s var(--n-bezier);
 display: inline-flex;
 align-items: center;
 `,[v("icon",`
 font-size: 18px;
 vertical-align: -.2em;
 transition: color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 `),k("&:not(:last-child)",[N("clickable",[m("link",`
 cursor: pointer;
 `,[k("&:hover",`
 background-color: var(--n-item-color-hover);
 `),k("&:active",`
 background-color: var(--n-item-color-pressed); 
 `)])])]),m("link",`
 padding: 4px;
 border-radius: var(--n-item-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 position: relative;
 `,[k("&:hover",`
 color: var(--n-item-text-color-hover);
 `,[v("icon",`
 color: var(--n-item-text-color-hover);
 `)]),k("&:active",`
 color: var(--n-item-text-color-pressed);
 `,[v("icon",`
 color: var(--n-item-text-color-pressed);
 `)])]),m("separator",`
 margin: 0 8px;
 color: var(--n-separator-color);
 transition: color .3s var(--n-bezier);
 user-select: none;
 -webkit-user-select: none;
 `),k("&:last-child",[m("link",`
 font-weight: var(--n-font-weight-active);
 cursor: unset;
 color: var(--n-item-text-color-active);
 `,[v("icon",`
 color: var(--n-item-text-color-active);
 `)]),m("separator",`
 display: none;
 `)])])]),Ge=ee("n-breadcrumb"),et=Object.assign(Object.assign({},q.props),{separator:{type:String,default:"/"}}),ot=S({name:"Breadcrumb",props:et,setup(e){const{mergedClsPrefixRef:o,inlineThemeDisabled:t}=te(e),i=q("Breadcrumb","-breadcrumb",Qo,go,e,o);Q(Ge,{separatorRef:ae(e,"separator"),mergedClsPrefixRef:o});const s=x(()=>{const{common:{cubicBezierEaseInOut:u},self:{separatorColor:c,itemTextColor:l,itemTextColorHover:h,itemTextColorPressed:M,itemTextColorActive:P,fontSize:p,fontWeightActive:y,itemBorderRadius:C,itemColorHover:f,itemColorPressed:I,itemLineHeight:A}}=i.value;return{"--n-font-size":p,"--n-bezier":u,"--n-item-text-color":l,"--n-item-text-color-hover":h,"--n-item-text-color-pressed":M,"--n-item-text-color-active":P,"--n-separator-color":c,"--n-item-color-hover":f,"--n-item-color-pressed":I,"--n-item-border-radius":C,"--n-font-weight-active":y,"--n-item-line-height":A}}),n=t?re("breadcrumb",void 0,s,e):void 0;return{mergedClsPrefix:o,cssVars:t?void 0:s,themeClass:n==null?void 0:n.themeClass,onRender:n==null?void 0:n.onRender}},render(){var e;return(e=this.onRender)===null||e===void 0||e.call(this),d("nav",{class:[`${this.mergedClsPrefix}-breadcrumb`,this.themeClass],style:this.cssVars,"aria-label":"Breadcrumb"},d("ul",null,this.$slots))}});function tt(e=jo?window:null){const o=()=>{const{hash:s,host:n,hostname:u,href:c,origin:l,pathname:h,port:M,protocol:P,search:p}=(e==null?void 0:e.location)||{};return{hash:s,host:n,hostname:u,href:c,origin:l,pathname:h,port:M,protocol:P,search:p}},t=F(o()),i=()=>{t.value=o()};return je(()=>{e&&(e.addEventListener("popstate",i),e.addEventListener("hashchange",i))}),bo(()=>{e&&(e.removeEventListener("popstate",i),e.removeEventListener("hashchange",i))}),t}const rt={separator:String,href:String,clickable:{type:Boolean,default:!0},showSeparator:{type:Boolean,default:!0},onClick:Function},Te=S({name:"BreadcrumbItem",props:rt,slots:Object,setup(e,{slots:o}){const t=Y(Ge,null);if(!t)return()=>null;const{separatorRef:i,mergedClsPrefixRef:s}=t,n=tt(),u=x(()=>e.href?"a":"span"),c=x(()=>n.value.href===e.href?"location":null);return()=>{const{value:l}=s;return d("li",{class:[`${l}-breadcrumb-item`,e.clickable&&`${l}-breadcrumb-item--clickable`]},d(u.value,{class:`${l}-breadcrumb-item__link`,"aria-current":c.value,href:e.href,onClick:e.onClick},o),e.showSeparator&&d("span",{class:`${l}-breadcrumb-item__separator`,"aria-hidden":"true"},Ee(o.separator,()=>{var h;return[(h=e.separator)!==null&&h!==void 0?h:i.value]})))}}});function nt(e){const{baseColor:o,textColor2:t,bodyColor:i,cardColor:s,dividerColor:n,actionColor:u,scrollbarColor:c,scrollbarColorHover:l,invertedColor:h}=e;return{textColor:t,textColorInverted:"#FFF",color:i,colorEmbedded:u,headerColor:s,headerColorInverted:h,footerColor:u,footerColorInverted:h,headerBorderColor:n,headerBorderColorInverted:h,footerBorderColor:n,footerBorderColorInverted:h,siderBorderColor:n,siderBorderColorInverted:h,siderColor:s,siderColorInverted:h,siderToggleButtonBorder:`1px solid ${n}`,siderToggleButtonColor:o,siderToggleButtonIconColor:t,siderToggleButtonIconColorInverted:t,siderToggleBarColor:Oe(i,c),siderToggleBarColorHover:Oe(i,l),__invertScrollbar:"true"}}const ze=xo({name:"Layout",common:wo,peers:{Scrollbar:Co},self:nt}),qe=ee("n-layout-sider"),Se={type:String,default:"static"},it=v("layout",`
 color: var(--n-text-color);
 background-color: var(--n-color);
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 flex: auto;
 overflow: hidden;
 transition:
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
`,[v("layout-scroll-container",`
 overflow-x: hidden;
 box-sizing: border-box;
 height: 100%;
 `),N("absolute-positioned",`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),lt={embedded:Boolean,position:Se,nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,onScroll:Function,contentClass:String,contentStyle:{type:[String,Object],default:""},hasSider:Boolean,siderPlacement:{type:String,default:"left"}},We=ee("n-layout");function Ye(e){return S({name:e?"LayoutContent":"Layout",props:Object.assign(Object.assign({},q.props),lt),setup(o){const t=F(null),i=F(null),{mergedClsPrefixRef:s,inlineThemeDisabled:n}=te(o),u=q("Layout","-layout",it,ze,o,s);function c(f,I){if(o.nativeScrollbar){const{value:A}=t;A&&(I===void 0?A.scrollTo(f):A.scrollTo(f,I))}else{const{value:A}=i;A&&A.scrollTo(f,I)}}Q(We,o);let l=0,h=0;const M=f=>{var I;const A=f.target;l=A.scrollLeft,h=A.scrollTop,(I=o.onScroll)===null||I===void 0||I.call(o,f)};Ve(()=>{if(o.nativeScrollbar){const f=t.value;f&&(f.scrollTop=h,f.scrollLeft=l)}});const P={display:"flex",flexWrap:"nowrap",width:"100%",flexDirection:"row"},p={scrollTo:c},y=x(()=>{const{common:{cubicBezierEaseInOut:f},self:I}=u.value;return{"--n-bezier":f,"--n-color":o.embedded?I.colorEmbedded:I.color,"--n-text-color":I.textColor}}),C=n?re("layout",x(()=>o.embedded?"e":""),y,o):void 0;return Object.assign({mergedClsPrefix:s,scrollableElRef:t,scrollbarInstRef:i,hasSiderStyle:P,mergedTheme:u,handleNativeElScroll:M,cssVars:n?void 0:y,themeClass:C==null?void 0:C.themeClass,onRender:C==null?void 0:C.onRender},p)},render(){var o;const{mergedClsPrefix:t,hasSider:i}=this;(o=this.onRender)===null||o===void 0||o.call(this);const s=i?this.hasSiderStyle:void 0,n=[this.themeClass,e&&`${t}-layout-content`,`${t}-layout`,`${t}-layout--${this.position}-positioned`];return d("div",{class:n,style:this.cssVars},this.nativeScrollbar?d("div",{ref:"scrollableElRef",class:[`${t}-layout-scroll-container`,this.contentClass],style:[this.contentStyle,s],onScroll:this.handleNativeElScroll},this.$slots):d(Ke,Object.assign({},this.scrollbarProps,{onScroll:this.onScroll,ref:"scrollbarInstRef",theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,contentClass:this.contentClass,contentStyle:[this.contentStyle,s]}),this.$slots))}})}const Me=Ye(!1),at=Ye(!0),st=v("layout-header",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 box-sizing: border-box;
 width: 100%;
 background-color: var(--n-color);
 color: var(--n-text-color);
`,[N("absolute-positioned",`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 `),N("bordered",`
 border-bottom: solid 1px var(--n-border-color);
 `)]),ct={position:Se,inverted:Boolean,bordered:{type:Boolean,default:!1}},dt=S({name:"LayoutHeader",props:Object.assign(Object.assign({},q.props),ct),setup(e){const{mergedClsPrefixRef:o,inlineThemeDisabled:t}=te(e),i=q("Layout","-layout-header",st,ze,e,o),s=x(()=>{const{common:{cubicBezierEaseInOut:u},self:c}=i.value,l={"--n-bezier":u};return e.inverted?(l["--n-color"]=c.headerColorInverted,l["--n-text-color"]=c.textColorInverted,l["--n-border-color"]=c.headerBorderColorInverted):(l["--n-color"]=c.headerColor,l["--n-text-color"]=c.textColor,l["--n-border-color"]=c.headerBorderColor),l}),n=t?re("layout-header",x(()=>e.inverted?"a":"b"),s,e):void 0;return{mergedClsPrefix:o,cssVars:t?void 0:s,themeClass:n==null?void 0:n.themeClass,onRender:n==null?void 0:n.onRender}},render(){var e;const{mergedClsPrefix:o}=this;return(e=this.onRender)===null||e===void 0||e.call(this),d("div",{class:[`${o}-layout-header`,this.themeClass,this.position&&`${o}-layout-header--${this.position}-positioned`,this.bordered&&`${o}-layout-header--bordered`],style:this.cssVars},this.$slots)}}),ut=v("layout-sider",`
 flex-shrink: 0;
 box-sizing: border-box;
 position: relative;
 z-index: 1;
 color: var(--n-text-color);
 transition:
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 min-width .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 transform .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 display: flex;
 justify-content: flex-end;
`,[N("bordered",[m("border",`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 width: 1px;
 background-color: var(--n-border-color);
 transition: background-color .3s var(--n-bezier);
 `)]),m("left-placement",[N("bordered",[m("border",`
 right: 0;
 `)])]),N("right-placement",`
 justify-content: flex-start;
 `,[N("bordered",[m("border",`
 left: 0;
 `)]),N("collapsed",[v("layout-toggle-button",[v("base-icon",`
 transform: rotate(180deg);
 `)]),v("layout-toggle-bar",[k("&:hover",[m("top",{transform:"rotate(-12deg) scale(1.15) translateY(-2px)"}),m("bottom",{transform:"rotate(12deg) scale(1.15) translateY(2px)"})])])]),v("layout-toggle-button",`
 left: 0;
 transform: translateX(-50%) translateY(-50%);
 `,[v("base-icon",`
 transform: rotate(0);
 `)]),v("layout-toggle-bar",`
 left: -28px;
 transform: rotate(180deg);
 `,[k("&:hover",[m("top",{transform:"rotate(12deg) scale(1.15) translateY(-2px)"}),m("bottom",{transform:"rotate(-12deg) scale(1.15) translateY(2px)"})])])]),N("collapsed",[v("layout-toggle-bar",[k("&:hover",[m("top",{transform:"rotate(-12deg) scale(1.15) translateY(-2px)"}),m("bottom",{transform:"rotate(12deg) scale(1.15) translateY(2px)"})])]),v("layout-toggle-button",[v("base-icon",`
 transform: rotate(0);
 `)])]),v("layout-toggle-button",`
 transition:
 color .3s var(--n-bezier),
 right .3s var(--n-bezier),
 left .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 cursor: pointer;
 width: 24px;
 height: 24px;
 position: absolute;
 top: 50%;
 right: 0;
 border-radius: 50%;
 display: flex;
 align-items: center;
 justify-content: center;
 font-size: 18px;
 color: var(--n-toggle-button-icon-color);
 border: var(--n-toggle-button-border);
 background-color: var(--n-toggle-button-color);
 box-shadow: 0 2px 4px 0px rgba(0, 0, 0, .06);
 transform: translateX(50%) translateY(-50%);
 z-index: 1;
 `,[v("base-icon",`
 transition: transform .3s var(--n-bezier);
 transform: rotate(180deg);
 `)]),v("layout-toggle-bar",`
 cursor: pointer;
 height: 72px;
 width: 32px;
 position: absolute;
 top: calc(50% - 36px);
 right: -28px;
 `,[m("top, bottom",`
 position: absolute;
 width: 4px;
 border-radius: 2px;
 height: 38px;
 left: 14px;
 transition: 
 background-color .3s var(--n-bezier),
 transform .3s var(--n-bezier);
 `),m("bottom",`
 position: absolute;
 top: 34px;
 `),k("&:hover",[m("top",{transform:"rotate(12deg) scale(1.15) translateY(-2px)"}),m("bottom",{transform:"rotate(-12deg) scale(1.15) translateY(2px)"})]),m("top, bottom",{backgroundColor:"var(--n-toggle-bar-color)"}),k("&:hover",[m("top, bottom",{backgroundColor:"var(--n-toggle-bar-color-hover)"})])]),m("border",`
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 width: 1px;
 transition: background-color .3s var(--n-bezier);
 `),v("layout-sider-scroll-container",`
 flex-grow: 1;
 flex-shrink: 0;
 box-sizing: border-box;
 height: 100%;
 opacity: 0;
 transition: opacity .3s var(--n-bezier);
 max-width: 100%;
 `),N("show-content",[v("layout-sider-scroll-container",{opacity:1})]),N("absolute-positioned",`
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 `)]),ht=S({props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){const{clsPrefix:e}=this;return d("div",{onClick:this.onClick,class:`${e}-layout-toggle-bar`},d("div",{class:`${e}-layout-toggle-bar__top`}),d("div",{class:`${e}-layout-toggle-bar__bottom`}))}}),vt=S({name:"LayoutToggleButton",props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){const{clsPrefix:e}=this;return d("div",{class:`${e}-layout-toggle-button`,onClick:this.onClick},d(De,{clsPrefix:e},{default:()=>d(Ho,null)}))}}),mt={position:Se,bordered:Boolean,collapsedWidth:{type:Number,default:48},width:{type:[Number,String],default:272},contentClass:String,contentStyle:{type:[String,Object],default:""},collapseMode:{type:String,default:"transform"},collapsed:{type:Boolean,default:void 0},defaultCollapsed:Boolean,showCollapsedContent:{type:Boolean,default:!0},showTrigger:{type:[Boolean,String],default:!1},nativeScrollbar:{type:Boolean,default:!0},inverted:Boolean,scrollbarProps:Object,triggerClass:String,triggerStyle:[String,Object],collapsedTriggerClass:String,collapsedTriggerStyle:[String,Object],"onUpdate:collapsed":[Function,Array],onUpdateCollapsed:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,onExpand:[Function,Array],onCollapse:[Function,Array],onScroll:Function},pt=S({name:"LayoutSider",props:Object.assign(Object.assign({},q.props),mt),setup(e){const o=Y(We),t=F(null),i=F(null),s=F(e.defaultCollapsed),n=xe(ae(e,"collapsed"),s),u=x(()=>me(n.value?e.collapsedWidth:e.width)),c=x(()=>e.collapseMode!=="transform"?{}:{minWidth:me(e.width)}),l=x(()=>o?o.siderPlacement:"left");function h(T,w){if(e.nativeScrollbar){const{value:R}=t;R&&(w===void 0?R.scrollTo(T):R.scrollTo(T,w))}else{const{value:R}=i;R&&R.scrollTo(T,w)}}function M(){const{"onUpdate:collapsed":T,onUpdateCollapsed:w,onExpand:R,onCollapse:U}=e,{value:G}=n;w&&W(w,!G),T&&W(T,!G),s.value=!G,G?R&&W(R):U&&W(U)}let P=0,p=0;const y=T=>{var w;const R=T.target;P=R.scrollLeft,p=R.scrollTop,(w=e.onScroll)===null||w===void 0||w.call(e,T)};Ve(()=>{if(e.nativeScrollbar){const T=t.value;T&&(T.scrollTop=p,T.scrollLeft=P)}}),Q(qe,{collapsedRef:n,collapseModeRef:ae(e,"collapseMode")});const{mergedClsPrefixRef:C,inlineThemeDisabled:f}=te(e),I=q("Layout","-layout-sider",ut,ze,e,C);function A(T){var w,R;T.propertyName==="max-width"&&(n.value?(w=e.onAfterLeave)===null||w===void 0||w.call(e):(R=e.onAfterEnter)===null||R===void 0||R.call(e))}const g={scrollTo:h},O=x(()=>{const{common:{cubicBezierEaseInOut:T},self:w}=I.value,{siderToggleButtonColor:R,siderToggleButtonBorder:U,siderToggleBarColor:G,siderToggleBarColorHover:ne}=w,V={"--n-bezier":T,"--n-toggle-button-color":R,"--n-toggle-button-border":U,"--n-toggle-bar-color":G,"--n-toggle-bar-color-hover":ne};return e.inverted?(V["--n-color"]=w.siderColorInverted,V["--n-text-color"]=w.textColorInverted,V["--n-border-color"]=w.siderBorderColorInverted,V["--n-toggle-button-icon-color"]=w.siderToggleButtonIconColorInverted,V.__invertScrollbar=w.__invertScrollbar):(V["--n-color"]=w.siderColor,V["--n-text-color"]=w.textColor,V["--n-border-color"]=w.siderBorderColor,V["--n-toggle-button-icon-color"]=w.siderToggleButtonIconColor),V}),H=f?re("layout-sider",x(()=>e.inverted?"a":"b"),O,e):void 0;return Object.assign({scrollableElRef:t,scrollbarInstRef:i,mergedClsPrefix:C,mergedTheme:I,styleMaxWidth:u,mergedCollapsed:n,scrollContainerStyle:c,siderPlacement:l,handleNativeElScroll:y,handleTransitionend:A,handleTriggerClick:M,inlineThemeDisabled:f,cssVars:O,themeClass:H==null?void 0:H.themeClass,onRender:H==null?void 0:H.onRender},g)},render(){var e;const{mergedClsPrefix:o,mergedCollapsed:t,showTrigger:i}=this;return(e=this.onRender)===null||e===void 0||e.call(this),d("aside",{class:[`${o}-layout-sider`,this.themeClass,`${o}-layout-sider--${this.position}-positioned`,`${o}-layout-sider--${this.siderPlacement}-placement`,this.bordered&&`${o}-layout-sider--bordered`,t&&`${o}-layout-sider--collapsed`,(!t||this.showCollapsedContent)&&`${o}-layout-sider--show-content`],onTransitionend:this.handleTransitionend,style:[this.inlineThemeDisabled?void 0:this.cssVars,{maxWidth:this.styleMaxWidth,width:me(this.width)}]},this.nativeScrollbar?d("div",{class:[`${o}-layout-sider-scroll-container`,this.contentClass],onScroll:this.handleNativeElScroll,style:[this.scrollContainerStyle,{overflow:"auto"},this.contentStyle],ref:"scrollableElRef"},this.$slots):d(Ke,Object.assign({},this.scrollbarProps,{onScroll:this.onScroll,ref:"scrollbarInstRef",style:this.scrollContainerStyle,contentStyle:this.contentStyle,contentClass:this.contentClass,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,builtinThemeOverrides:this.inverted&&this.cssVars.__invertScrollbar==="true"?{colorHover:"rgba(255, 255, 255, .4)",color:"rgba(255, 255, 255, .3)"}:void 0}),this.$slots),i?i==="bar"?d(ht,{clsPrefix:o,class:t?this.collapsedTriggerClass:this.triggerClass,style:t?this.collapsedTriggerStyle:this.triggerStyle,onClick:this.handleTriggerClick}):d(vt,{clsPrefix:o,class:t?this.collapsedTriggerClass:this.triggerClass,style:t?this.collapsedTriggerStyle:this.triggerStyle,onClick:this.handleTriggerClick}):null,this.bordered?d("div",{class:`${o}-layout-sider__border`}):null)}}),se=ee("n-menu"),Xe=ee("n-submenu"),Ie=ee("n-menu-item-group"),Ne=[k("&::before","background-color: var(--n-item-color-hover);"),m("arrow",`
 color: var(--n-arrow-color-hover);
 `),m("icon",`
 color: var(--n-item-icon-color-hover);
 `),v("menu-item-content-header",`
 color: var(--n-item-text-color-hover);
 `,[k("a",`
 color: var(--n-item-text-color-hover);
 `),m("extra",`
 color: var(--n-item-text-color-hover);
 `)])],Be=[m("icon",`
 color: var(--n-item-icon-color-hover-horizontal);
 `),v("menu-item-content-header",`
 color: var(--n-item-text-color-hover-horizontal);
 `,[k("a",`
 color: var(--n-item-text-color-hover-horizontal);
 `),m("extra",`
 color: var(--n-item-text-color-hover-horizontal);
 `)])],ft=k([v("menu",`
 background-color: var(--n-color);
 color: var(--n-item-text-color);
 overflow: hidden;
 transition: background-color .3s var(--n-bezier);
 box-sizing: border-box;
 font-size: var(--n-font-size);
 padding-bottom: 6px;
 `,[N("horizontal",`
 max-width: 100%;
 width: 100%;
 display: flex;
 overflow: hidden;
 padding-bottom: 0;
 `,[v("submenu","margin: 0;"),v("menu-item","margin: 0;"),v("menu-item-content",`
 padding: 0 20px;
 border-bottom: 2px solid #0000;
 `,[k("&::before","display: none;"),N("selected","border-bottom: 2px solid var(--n-border-color-horizontal)")]),v("menu-item-content",[N("selected",[m("icon","color: var(--n-item-icon-color-active-horizontal);"),v("menu-item-content-header",`
 color: var(--n-item-text-color-active-horizontal);
 `,[k("a","color: var(--n-item-text-color-active-horizontal);"),m("extra","color: var(--n-item-text-color-active-horizontal);")])]),N("child-active",`
 border-bottom: 2px solid var(--n-border-color-horizontal);
 `,[v("menu-item-content-header",`
 color: var(--n-item-text-color-child-active-horizontal);
 `,[k("a",`
 color: var(--n-item-text-color-child-active-horizontal);
 `),m("extra",`
 color: var(--n-item-text-color-child-active-horizontal);
 `)]),m("icon",`
 color: var(--n-item-icon-color-child-active-horizontal);
 `)]),le("disabled",[le("selected, child-active",[k("&:focus-within",Be)]),N("selected",[X(null,[m("icon","color: var(--n-item-icon-color-active-hover-horizontal);"),v("menu-item-content-header",`
 color: var(--n-item-text-color-active-hover-horizontal);
 `,[k("a","color: var(--n-item-text-color-active-hover-horizontal);"),m("extra","color: var(--n-item-text-color-active-hover-horizontal);")])])]),N("child-active",[X(null,[m("icon","color: var(--n-item-icon-color-child-active-hover-horizontal);"),v("menu-item-content-header",`
 color: var(--n-item-text-color-child-active-hover-horizontal);
 `,[k("a","color: var(--n-item-text-color-child-active-hover-horizontal);"),m("extra","color: var(--n-item-text-color-child-active-hover-horizontal);")])])]),X("border-bottom: 2px solid var(--n-border-color-horizontal);",Be)]),v("menu-item-content-header",[k("a","color: var(--n-item-text-color-horizontal);")])])]),le("responsive",[v("menu-item-content-header",`
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),N("collapsed",[v("menu-item-content",[N("selected",[k("&::before",`
 background-color: var(--n-item-color-active-collapsed) !important;
 `)]),v("menu-item-content-header","opacity: 0;"),m("arrow","opacity: 0;"),m("icon","color: var(--n-item-icon-color-collapsed);")])]),v("menu-item",`
 height: var(--n-item-height);
 margin-top: 6px;
 position: relative;
 `),v("menu-item-content",`
 box-sizing: border-box;
 line-height: 1.75;
 height: 100%;
 display: grid;
 grid-template-areas: "icon content arrow";
 grid-template-columns: auto 1fr auto;
 align-items: center;
 cursor: pointer;
 position: relative;
 padding-right: 18px;
 transition:
 background-color .3s var(--n-bezier),
 padding-left .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[k("> *","z-index: 1;"),k("&::before",`
 z-index: auto;
 content: "";
 background-color: #0000;
 position: absolute;
 left: 8px;
 right: 8px;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),N("disabled",`
 opacity: .45;
 cursor: not-allowed;
 `),N("collapsed",[m("arrow","transform: rotate(0);")]),N("selected",[k("&::before","background-color: var(--n-item-color-active);"),m("arrow","color: var(--n-arrow-color-active);"),m("icon","color: var(--n-item-icon-color-active);"),v("menu-item-content-header",`
 color: var(--n-item-text-color-active);
 `,[k("a","color: var(--n-item-text-color-active);"),m("extra","color: var(--n-item-text-color-active);")])]),N("child-active",[v("menu-item-content-header",`
 color: var(--n-item-text-color-child-active);
 `,[k("a",`
 color: var(--n-item-text-color-child-active);
 `),m("extra",`
 color: var(--n-item-text-color-child-active);
 `)]),m("arrow",`
 color: var(--n-arrow-color-child-active);
 `),m("icon",`
 color: var(--n-item-icon-color-child-active);
 `)]),le("disabled",[le("selected, child-active",[k("&:focus-within",Ne)]),N("selected",[X(null,[m("arrow","color: var(--n-arrow-color-active-hover);"),m("icon","color: var(--n-item-icon-color-active-hover);"),v("menu-item-content-header",`
 color: var(--n-item-text-color-active-hover);
 `,[k("a","color: var(--n-item-text-color-active-hover);"),m("extra","color: var(--n-item-text-color-active-hover);")])])]),N("child-active",[X(null,[m("arrow","color: var(--n-arrow-color-child-active-hover);"),m("icon","color: var(--n-item-icon-color-child-active-hover);"),v("menu-item-content-header",`
 color: var(--n-item-text-color-child-active-hover);
 `,[k("a","color: var(--n-item-text-color-child-active-hover);"),m("extra","color: var(--n-item-text-color-child-active-hover);")])])]),N("selected",[X(null,[k("&::before","background-color: var(--n-item-color-active-hover);")])]),X(null,Ne)]),m("icon",`
 grid-area: icon;
 color: var(--n-item-icon-color);
 transition:
 color .3s var(--n-bezier),
 font-size .3s var(--n-bezier),
 margin-right .3s var(--n-bezier);
 box-sizing: content-box;
 display: inline-flex;
 align-items: center;
 justify-content: center;
 `),m("arrow",`
 grid-area: arrow;
 font-size: 16px;
 color: var(--n-arrow-color);
 transform: rotate(180deg);
 opacity: 1;
 transition:
 color .3s var(--n-bezier),
 transform 0.2s var(--n-bezier),
 opacity 0.2s var(--n-bezier);
 `),v("menu-item-content-header",`
 grid-area: content;
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 opacity: 1;
 white-space: nowrap;
 color: var(--n-item-text-color);
 `,[k("a",`
 outline: none;
 text-decoration: none;
 transition: color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 `,[k("&::before",`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),m("extra",`
 font-size: .93em;
 color: var(--n-group-text-color);
 transition: color .3s var(--n-bezier);
 `)])]),v("submenu",`
 cursor: pointer;
 position: relative;
 margin-top: 6px;
 `,[v("menu-item-content",`
 height: var(--n-item-height);
 `),v("submenu-children",`
 overflow: hidden;
 padding: 0;
 `,[To({duration:".2s"})])]),v("menu-item-group",[v("menu-item-group-title",`
 margin-top: 6px;
 color: var(--n-group-text-color);
 cursor: default;
 font-size: .93em;
 height: 36px;
 display: flex;
 align-items: center;
 transition:
 padding-left .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)])]),v("menu-tooltip",[k("a",`
 color: inherit;
 text-decoration: none;
 `)]),v("menu-divider",`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 6px 18px;
 `)]);function X(e,o){return[N("hover",e,o),k("&:hover",e,o)]}const Ze=S({name:"MenuOptionContent",props:{collapsed:Boolean,disabled:Boolean,title:[String,Function],icon:Function,extra:[String,Function],showArrow:Boolean,childActive:Boolean,hover:Boolean,paddingLeft:Number,selected:Boolean,maxIconSize:{type:Number,required:!0},activeIconSize:{type:Number,required:!0},iconMarginRight:{type:Number,required:!0},clsPrefix:{type:String,required:!0},onClick:Function,tmNode:{type:Object,required:!0},isEllipsisPlaceholder:Boolean},setup(e){const{props:o}=Y(se);return{menuProps:o,style:x(()=>{const{paddingLeft:t}=e;return{paddingLeft:t&&`${t}px`}}),iconStyle:x(()=>{const{maxIconSize:t,activeIconSize:i,iconMarginRight:s}=e;return{width:`${t}px`,height:`${t}px`,fontSize:`${i}px`,marginRight:`${s}px`}})}},render(){const{clsPrefix:e,tmNode:o,menuProps:{renderIcon:t,renderLabel:i,renderExtra:s,expandIcon:n}}=this,u=t?t(o.rawNode):oe(this.icon);return d("div",{onClick:c=>{var l;(l=this.onClick)===null||l===void 0||l.call(this,c)},role:"none",class:[`${e}-menu-item-content`,{[`${e}-menu-item-content--selected`]:this.selected,[`${e}-menu-item-content--collapsed`]:this.collapsed,[`${e}-menu-item-content--child-active`]:this.childActive,[`${e}-menu-item-content--disabled`]:this.disabled,[`${e}-menu-item-content--hover`]:this.hover}],style:this.style},u&&d("div",{class:`${e}-menu-item-content__icon`,style:this.iconStyle,role:"none"},[u]),d("div",{class:`${e}-menu-item-content-header`,role:"none"},this.isEllipsisPlaceholder?this.title:i?i(o.rawNode):oe(this.title),this.extra||s?d("span",{class:`${e}-menu-item-content-header__extra`}," ",s?s(o.rawNode):oe(this.extra)):null),this.showArrow?d(De,{ariaHidden:!0,class:`${e}-menu-item-content__arrow`,clsPrefix:e},{default:()=>n?n(o.rawNode):d(Wo,null)}):null)}}),ue=8;function Re(e){const o=Y(se),{props:t,mergedCollapsedRef:i}=o,s=Y(Xe,null),n=Y(Ie,null),u=x(()=>t.mode==="horizontal"),c=x(()=>u.value?t.dropdownPlacement:"tmNodes"in e?"right-start":"right"),l=x(()=>{var p;return Math.max((p=t.collapsedIconSize)!==null&&p!==void 0?p:t.iconSize,t.iconSize)}),h=x(()=>{var p;return!u.value&&e.root&&i.value&&(p=t.collapsedIconSize)!==null&&p!==void 0?p:t.iconSize}),M=x(()=>{if(u.value)return;const{collapsedWidth:p,indent:y,rootIndent:C}=t,{root:f,isGroup:I}=e,A=C===void 0?y:C;return f?i.value?p/2-l.value/2:A:n&&typeof n.paddingLeftRef.value=="number"?y/2+n.paddingLeftRef.value:s&&typeof s.paddingLeftRef.value=="number"?(I?y/2:y)+s.paddingLeftRef.value:0}),P=x(()=>{const{collapsedWidth:p,indent:y,rootIndent:C}=t,{value:f}=l,{root:I}=e;return u.value||!I||!i.value?ue:(C===void 0?y:C)+f+ue-(p+f)/2});return{dropdownPlacement:c,activeIconSize:h,maxIconSize:l,paddingLeft:M,iconMarginRight:P,NMenu:o,NSubmenu:s,NMenuOptionGroup:n}}const $e={internalKey:{type:[String,Number],required:!0},root:Boolean,isGroup:Boolean,level:{type:Number,required:!0},title:[String,Function],extra:[String,Function]},gt=S({name:"MenuDivider",setup(){const e=Y(se),{mergedClsPrefixRef:o,isHorizontalRef:t}=e;return()=>t.value?null:d("div",{class:`${o.value}-menu-divider`})}}),Je=Object.assign(Object.assign({},$e),{tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function}),bt=_e(Je),xt=S({name:"MenuOption",props:Je,setup(e){const o=Re(e),{NSubmenu:t,NMenu:i,NMenuOptionGroup:s}=o,{props:n,mergedClsPrefixRef:u,mergedCollapsedRef:c}=i,l=t?t.mergedDisabledRef:s?s.mergedDisabledRef:{value:!1},h=x(()=>l.value||e.disabled);function M(p){const{onClick:y}=e;y&&y(p)}function P(p){h.value||(i.doSelect(e.internalKey,e.tmNode.rawNode),M(p))}return{mergedClsPrefix:u,dropdownPlacement:o.dropdownPlacement,paddingLeft:o.paddingLeft,iconMarginRight:o.iconMarginRight,maxIconSize:o.maxIconSize,activeIconSize:o.activeIconSize,mergedTheme:i.mergedThemeRef,menuProps:n,dropdownEnabled:be(()=>e.root&&c.value&&n.mode!=="horizontal"&&!h.value),selected:be(()=>i.mergedValueRef.value===e.internalKey),mergedDisabled:h,handleClick:P}},render(){const{mergedClsPrefix:e,mergedTheme:o,tmNode:t,menuProps:{renderLabel:i,nodeProps:s}}=this,n=s==null?void 0:s(t.rawNode);return d("div",Object.assign({},n,{role:"menuitem",class:[`${e}-menu-item`,n==null?void 0:n.class]}),d(Ao,{theme:o.peers.Tooltip,themeOverrides:o.peerOverrides.Tooltip,trigger:"hover",placement:this.dropdownPlacement,disabled:!this.dropdownEnabled||this.title===void 0,internalExtraClass:["menu-tooltip"]},{default:()=>i?i(t.rawNode):oe(this.title),trigger:()=>d(Ze,{tmNode:t,clsPrefix:e,paddingLeft:this.paddingLeft,iconMarginRight:this.iconMarginRight,maxIconSize:this.maxIconSize,activeIconSize:this.activeIconSize,selected:this.selected,title:this.title,extra:this.extra,disabled:this.mergedDisabled,icon:this.icon,onClick:this.handleClick})}))}}),Qe=Object.assign(Object.assign({},$e),{tmNode:{type:Object,required:!0},tmNodes:{type:Array,required:!0}}),Ct=_e(Qe),wt=S({name:"MenuOptionGroup",props:Qe,setup(e){const o=Re(e),{NSubmenu:t}=o,i=x(()=>t!=null&&t.mergedDisabledRef.value?!0:e.tmNode.disabled);Q(Ie,{paddingLeftRef:o.paddingLeft,mergedDisabledRef:i});const{mergedClsPrefixRef:s,props:n}=Y(se);return function(){const{value:u}=s,c=o.paddingLeft.value,{nodeProps:l}=n,h=l==null?void 0:l(e.tmNode.rawNode);return d("div",{class:`${u}-menu-item-group`,role:"group"},d("div",Object.assign({},h,{class:[`${u}-menu-item-group-title`,h==null?void 0:h.class],style:[(h==null?void 0:h.style)||"",c!==void 0?`padding-left: ${c}px;`:""]}),oe(e.title),e.extra?d(Ae,null," ",oe(e.extra)):null),d("div",null,e.tmNodes.map(M=>Pe(M,n))))}}});function Ce(e){return e.type==="divider"||e.type==="render"}function kt(e){return e.type==="divider"}function Pe(e,o){const{rawNode:t}=e,{show:i}=t;if(i===!1)return null;if(Ce(t))return kt(t)?d(gt,Object.assign({key:e.key},t.props)):null;const{labelField:s}=o,{key:n,level:u,isGroup:c}=e,l=Object.assign(Object.assign({},t),{title:t.title||t[s],extra:t.titleExtra||t.extra,key:n,internalKey:n,level:u,root:u===0,isGroup:c});return e.children?e.isGroup?d(wt,pe(l,Ct,{tmNode:e,tmNodes:e.children,key:n})):d(we,pe(l,yt,{key:n,rawNodes:t[o.childrenField],tmNodes:e.children,tmNode:e})):d(xt,pe(l,bt,{key:n,tmNode:e}))}const eo=Object.assign(Object.assign({},$e),{rawNodes:{type:Array,default:()=>[]},tmNodes:{type:Array,default:()=>[]},tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function,domId:String,virtualChildActive:{type:Boolean,default:void 0},isEllipsisPlaceholder:Boolean}),yt=_e(eo),we=S({name:"Submenu",props:eo,setup(e){const o=Re(e),{NMenu:t,NSubmenu:i}=o,{props:s,mergedCollapsedRef:n,mergedThemeRef:u}=t,c=x(()=>{const{disabled:p}=e;return i!=null&&i.mergedDisabledRef.value||s.disabled?!0:p}),l=F(!1);Q(Xe,{paddingLeftRef:o.paddingLeft,mergedDisabledRef:c}),Q(Ie,null);function h(){const{onClick:p}=e;p&&p()}function M(){c.value||(n.value||t.toggleExpand(e.internalKey),h())}function P(p){l.value=p}return{menuProps:s,mergedTheme:u,doSelect:t.doSelect,inverted:t.invertedRef,isHorizontal:t.isHorizontalRef,mergedClsPrefix:t.mergedClsPrefixRef,maxIconSize:o.maxIconSize,activeIconSize:o.activeIconSize,iconMarginRight:o.iconMarginRight,dropdownPlacement:o.dropdownPlacement,dropdownShow:l,paddingLeft:o.paddingLeft,mergedDisabled:c,mergedValue:t.mergedValueRef,childActive:be(()=>{var p;return(p=e.virtualChildActive)!==null&&p!==void 0?p:t.activePathRef.value.includes(e.internalKey)}),collapsed:x(()=>s.mode==="horizontal"?!1:n.value?!0:!t.mergedExpandedKeysRef.value.includes(e.internalKey)),dropdownEnabled:x(()=>!c.value&&(s.mode==="horizontal"||n.value)),handlePopoverShowChange:P,handleClick:M}},render(){var e;const{mergedClsPrefix:o,menuProps:{renderIcon:t,renderLabel:i}}=this,s=()=>{const{isHorizontal:u,paddingLeft:c,collapsed:l,mergedDisabled:h,maxIconSize:M,activeIconSize:P,title:p,childActive:y,icon:C,handleClick:f,menuProps:{nodeProps:I},dropdownShow:A,iconMarginRight:g,tmNode:O,mergedClsPrefix:H,isEllipsisPlaceholder:T,extra:w}=this,R=I==null?void 0:I(O.rawNode);return d("div",Object.assign({},R,{class:[`${H}-menu-item`,R==null?void 0:R.class],role:"menuitem"}),d(Ze,{tmNode:O,paddingLeft:c,collapsed:l,disabled:h,iconMarginRight:g,maxIconSize:M,activeIconSize:P,title:p,extra:w,showArrow:!u,childActive:y,clsPrefix:H,icon:C,hover:A,onClick:f,isEllipsisPlaceholder:T}))},n=()=>d(Lo,null,{default:()=>{const{tmNodes:u,collapsed:c}=this;return c?null:d("div",{class:`${o}-submenu-children`,role:"menu"},u.map(l=>Pe(l,this.menuProps)))}});return this.root?d(Ue,Object.assign({size:"large",trigger:"hover"},(e=this.menuProps)===null||e===void 0?void 0:e.dropdownProps,{themeOverrides:this.mergedTheme.peerOverrides.Dropdown,theme:this.mergedTheme.peers.Dropdown,builtinThemeOverrides:{fontSizeLarge:"14px",optionIconSizeLarge:"18px"},value:this.mergedValue,disabled:!this.dropdownEnabled,placement:this.dropdownPlacement,keyField:this.menuProps.keyField,labelField:this.menuProps.labelField,childrenField:this.menuProps.childrenField,onUpdateShow:this.handlePopoverShowChange,options:this.rawNodes,onSelect:this.doSelect,inverted:this.inverted,renderIcon:t,renderLabel:i}),{default:()=>d("div",{class:`${o}-submenu`,role:"menu","aria-expanded":!this.collapsed,id:this.domId},s(),this.isHorizontal?null:n())}):d("div",{class:`${o}-submenu`,role:"menu","aria-expanded":!this.collapsed,id:this.domId},s(),n())}}),_t=Object.assign(Object.assign({},q.props),{options:{type:Array,default:()=>[]},collapsed:{type:Boolean,default:void 0},collapsedWidth:{type:Number,default:48},iconSize:{type:Number,default:20},collapsedIconSize:{type:Number,default:24},rootIndent:Number,indent:{type:Number,default:32},labelField:{type:String,default:"label"},keyField:{type:String,default:"key"},childrenField:{type:String,default:"children"},disabledField:{type:String,default:"disabled"},defaultExpandAll:Boolean,defaultExpandedKeys:Array,expandedKeys:Array,value:[String,Number],defaultValue:{type:[String,Number],default:null},mode:{type:String,default:"vertical"},watchProps:{type:Array,default:void 0},disabled:Boolean,show:{type:Boolean,default:!0},inverted:Boolean,"onUpdate:expandedKeys":[Function,Array],onUpdateExpandedKeys:[Function,Array],onUpdateValue:[Function,Array],"onUpdate:value":[Function,Array],expandIcon:Function,renderIcon:Function,renderLabel:Function,renderExtra:Function,dropdownProps:Object,accordion:Boolean,nodeProps:Function,dropdownPlacement:{type:String,default:"bottom"},responsive:Boolean,items:Array,onOpenNamesChange:[Function,Array],onSelect:[Function,Array],onExpandedNamesChange:[Function,Array],expandedNames:Array,defaultExpandedNames:Array}),zt=S({name:"Menu",inheritAttrs:!1,props:_t,setup(e){const{mergedClsPrefixRef:o,inlineThemeDisabled:t}=te(e),i=q("Menu","-menu",ft,yo,e,o),s=Y(qe,null),n=x(()=>{var b;const{collapsed:$}=e;if($!==void 0)return $;if(s){const{collapseModeRef:r,collapsedRef:_}=s;if(r.value==="width")return(b=_.value)!==null&&b!==void 0?b:!1}return!1}),u=x(()=>{const{keyField:b,childrenField:$,disabledField:r}=e;return ve(e.items||e.options,{getIgnored(_){return Ce(_)},getChildren(_){return _[$]},getDisabled(_){return _[r]},getKey(_){var L;return(L=_[b])!==null&&L!==void 0?L:_.name}})}),c=x(()=>new Set(u.value.treeNodes.map(b=>b.key))),{watchProps:l}=e,h=F(null);l!=null&&l.includes("defaultValue")?ge(()=>{h.value=e.defaultValue}):h.value=e.defaultValue;const M=ae(e,"value"),P=xe(M,h),p=F([]),y=()=>{p.value=e.defaultExpandAll?u.value.getNonLeafKeys():e.defaultExpandedNames||e.defaultExpandedKeys||u.value.getPath(P.value,{includeSelf:!1}).keyPath};l!=null&&l.includes("defaultExpandedKeys")?ge(y):y();const C=Fo(e,["expandedNames","expandedKeys"]),f=xe(C,p),I=x(()=>u.value.treeNodes),A=x(()=>u.value.getPath(P.value).keyPath);Q(se,{props:e,mergedCollapsedRef:n,mergedThemeRef:i,mergedValueRef:P,mergedExpandedKeysRef:f,activePathRef:A,mergedClsPrefixRef:o,isHorizontalRef:x(()=>e.mode==="horizontal"),invertedRef:ae(e,"inverted"),doSelect:g,toggleExpand:H});function g(b,$){const{"onUpdate:value":r,onUpdateValue:_,onSelect:L}=e;_&&W(_,b,$),r&&W(r,b,$),L&&W(L,b,$),h.value=b}function O(b){const{"onUpdate:expandedKeys":$,onUpdateExpandedKeys:r,onExpandedNamesChange:_,onOpenNamesChange:L}=e;$&&W($,b),r&&W(r,b),_&&W(_,b),L&&W(L,b),p.value=b}function H(b){const $=Array.from(f.value),r=$.findIndex(_=>_===b);if(~r)$.splice(r,1);else{if(e.accordion&&c.value.has(b)){const _=$.findIndex(L=>c.value.has(L));_>-1&&$.splice(_,1)}$.push(b)}O($)}const T=b=>{const $=u.value.getPath(b??P.value,{includeSelf:!1}).keyPath;if(!$.length)return;const r=Array.from(f.value),_=new Set([...r,...$]);e.accordion&&c.value.forEach(L=>{_.has(L)&&!$.includes(L)&&_.delete(L)}),O(Array.from(_))},w=x(()=>{const{inverted:b}=e,{common:{cubicBezierEaseInOut:$},self:r}=i.value,{borderRadius:_,borderColorHorizontal:L,fontSize:ao,itemHeight:so,dividerColor:co}=r,a={"--n-divider-color":co,"--n-bezier":$,"--n-font-size":ao,"--n-border-color-horizontal":L,"--n-border-radius":_,"--n-item-height":so};return b?(a["--n-group-text-color"]=r.groupTextColorInverted,a["--n-color"]=r.colorInverted,a["--n-item-text-color"]=r.itemTextColorInverted,a["--n-item-text-color-hover"]=r.itemTextColorHoverInverted,a["--n-item-text-color-active"]=r.itemTextColorActiveInverted,a["--n-item-text-color-child-active"]=r.itemTextColorChildActiveInverted,a["--n-item-text-color-child-active-hover"]=r.itemTextColorChildActiveInverted,a["--n-item-text-color-active-hover"]=r.itemTextColorActiveHoverInverted,a["--n-item-icon-color"]=r.itemIconColorInverted,a["--n-item-icon-color-hover"]=r.itemIconColorHoverInverted,a["--n-item-icon-color-active"]=r.itemIconColorActiveInverted,a["--n-item-icon-color-active-hover"]=r.itemIconColorActiveHoverInverted,a["--n-item-icon-color-child-active"]=r.itemIconColorChildActiveInverted,a["--n-item-icon-color-child-active-hover"]=r.itemIconColorChildActiveHoverInverted,a["--n-item-icon-color-collapsed"]=r.itemIconColorCollapsedInverted,a["--n-item-text-color-horizontal"]=r.itemTextColorHorizontalInverted,a["--n-item-text-color-hover-horizontal"]=r.itemTextColorHoverHorizontalInverted,a["--n-item-text-color-active-horizontal"]=r.itemTextColorActiveHorizontalInverted,a["--n-item-text-color-child-active-horizontal"]=r.itemTextColorChildActiveHorizontalInverted,a["--n-item-text-color-child-active-hover-horizontal"]=r.itemTextColorChildActiveHoverHorizontalInverted,a["--n-item-text-color-active-hover-horizontal"]=r.itemTextColorActiveHoverHorizontalInverted,a["--n-item-icon-color-horizontal"]=r.itemIconColorHorizontalInverted,a["--n-item-icon-color-hover-horizontal"]=r.itemIconColorHoverHorizontalInverted,a["--n-item-icon-color-active-horizontal"]=r.itemIconColorActiveHorizontalInverted,a["--n-item-icon-color-active-hover-horizontal"]=r.itemIconColorActiveHoverHorizontalInverted,a["--n-item-icon-color-child-active-horizontal"]=r.itemIconColorChildActiveHorizontalInverted,a["--n-item-icon-color-child-active-hover-horizontal"]=r.itemIconColorChildActiveHoverHorizontalInverted,a["--n-arrow-color"]=r.arrowColorInverted,a["--n-arrow-color-hover"]=r.arrowColorHoverInverted,a["--n-arrow-color-active"]=r.arrowColorActiveInverted,a["--n-arrow-color-active-hover"]=r.arrowColorActiveHoverInverted,a["--n-arrow-color-child-active"]=r.arrowColorChildActiveInverted,a["--n-arrow-color-child-active-hover"]=r.arrowColorChildActiveHoverInverted,a["--n-item-color-hover"]=r.itemColorHoverInverted,a["--n-item-color-active"]=r.itemColorActiveInverted,a["--n-item-color-active-hover"]=r.itemColorActiveHoverInverted,a["--n-item-color-active-collapsed"]=r.itemColorActiveCollapsedInverted):(a["--n-group-text-color"]=r.groupTextColor,a["--n-color"]=r.color,a["--n-item-text-color"]=r.itemTextColor,a["--n-item-text-color-hover"]=r.itemTextColorHover,a["--n-item-text-color-active"]=r.itemTextColorActive,a["--n-item-text-color-child-active"]=r.itemTextColorChildActive,a["--n-item-text-color-child-active-hover"]=r.itemTextColorChildActiveHover,a["--n-item-text-color-active-hover"]=r.itemTextColorActiveHover,a["--n-item-icon-color"]=r.itemIconColor,a["--n-item-icon-color-hover"]=r.itemIconColorHover,a["--n-item-icon-color-active"]=r.itemIconColorActive,a["--n-item-icon-color-active-hover"]=r.itemIconColorActiveHover,a["--n-item-icon-color-child-active"]=r.itemIconColorChildActive,a["--n-item-icon-color-child-active-hover"]=r.itemIconColorChildActiveHover,a["--n-item-icon-color-collapsed"]=r.itemIconColorCollapsed,a["--n-item-text-color-horizontal"]=r.itemTextColorHorizontal,a["--n-item-text-color-hover-horizontal"]=r.itemTextColorHoverHorizontal,a["--n-item-text-color-active-horizontal"]=r.itemTextColorActiveHorizontal,a["--n-item-text-color-child-active-horizontal"]=r.itemTextColorChildActiveHorizontal,a["--n-item-text-color-child-active-hover-horizontal"]=r.itemTextColorChildActiveHoverHorizontal,a["--n-item-text-color-active-hover-horizontal"]=r.itemTextColorActiveHoverHorizontal,a["--n-item-icon-color-horizontal"]=r.itemIconColorHorizontal,a["--n-item-icon-color-hover-horizontal"]=r.itemIconColorHoverHorizontal,a["--n-item-icon-color-active-horizontal"]=r.itemIconColorActiveHorizontal,a["--n-item-icon-color-active-hover-horizontal"]=r.itemIconColorActiveHoverHorizontal,a["--n-item-icon-color-child-active-horizontal"]=r.itemIconColorChildActiveHorizontal,a["--n-item-icon-color-child-active-hover-horizontal"]=r.itemIconColorChildActiveHoverHorizontal,a["--n-arrow-color"]=r.arrowColor,a["--n-arrow-color-hover"]=r.arrowColorHover,a["--n-arrow-color-active"]=r.arrowColorActive,a["--n-arrow-color-active-hover"]=r.arrowColorActiveHover,a["--n-arrow-color-child-active"]=r.arrowColorChildActive,a["--n-arrow-color-child-active-hover"]=r.arrowColorChildActiveHover,a["--n-item-color-hover"]=r.itemColorHover,a["--n-item-color-active"]=r.itemColorActive,a["--n-item-color-active-hover"]=r.itemColorActiveHover,a["--n-item-color-active-collapsed"]=r.itemColorActiveCollapsed),a}),R=t?re("menu",x(()=>e.inverted?"a":"b"),w,e):void 0,U=Eo(),G=F(null),ne=F(null);let V=!0;const ce=()=>{var b;V?V=!1:(b=G.value)===null||b===void 0||b.sync({showAllItemsBeforeCalculate:!0})};function ie(){return document.getElementById(U)}const de=F(-1);function oo(b){de.value=e.options.length-b}function to(b){b||(de.value=-1)}const ro=x(()=>{const b=de.value;return{children:b===-1?[]:e.options.slice(b)}}),no=x(()=>{const{childrenField:b,disabledField:$,keyField:r}=e;return ve([ro.value],{getIgnored(_){return Ce(_)},getChildren(_){return _[b]},getDisabled(_){return _[$]},getKey(_){var L;return(L=_[r])!==null&&L!==void 0?L:_.name}})}),io=x(()=>ve([{}]).treeNodes[0]);function lo(){var b;if(de.value===-1)return d(we,{root:!0,level:0,key:"__ellpisisGroupPlaceholder__",internalKey:"__ellpisisGroupPlaceholder__",title:"···",tmNode:io.value,domId:U,isEllipsisPlaceholder:!0});const $=no.value.treeNodes[0],r=A.value,_=!!(!((b=$.children)===null||b===void 0)&&b.some(L=>r.includes(L.key)));return d(we,{level:0,root:!0,key:"__ellpisisGroup__",internalKey:"__ellpisisGroup__",title:"···",virtualChildActive:_,tmNode:$,domId:U,rawNodes:$.rawNode.children||[],tmNodes:$.children||[],isEllipsisPlaceholder:!0})}return{mergedClsPrefix:o,controlledExpandedKeys:C,uncontrolledExpanededKeys:p,mergedExpandedKeys:f,uncontrolledValue:h,mergedValue:P,activePath:A,tmNodes:I,mergedTheme:i,mergedCollapsed:n,cssVars:t?void 0:w,themeClass:R==null?void 0:R.themeClass,overflowRef:G,counterRef:ne,updateCounter:()=>{},onResize:ce,onUpdateOverflow:to,onUpdateCount:oo,renderCounter:lo,getCounter:ie,onRender:R==null?void 0:R.onRender,showOption:T,deriveResponsiveState:ce}},render(){const{mergedClsPrefix:e,mode:o,themeClass:t,onRender:i}=this;i==null||i();const s=()=>this.tmNodes.map(l=>Pe(l,this.$props)),u=o==="horizontal"&&this.responsive,c=()=>d("div",ko(this.$attrs,{role:o==="horizontal"?"menubar":"menu",class:[`${e}-menu`,t,`${e}-menu--${o}`,u&&`${e}-menu--responsive`,this.mergedCollapsed&&`${e}-menu--collapsed`],style:this.cssVars}),u?d(Mo,{ref:"overflowRef",onUpdateOverflow:this.onUpdateOverflow,getCounter:this.getCounter,onUpdateCount:this.onUpdateCount,updateCounter:this.updateCounter,style:{width:"100%",display:"flex",overflow:"hidden"}},{default:s,counter:this.renderCounter}):s());return u?d(Fe,{onResize:this.onResize},{default:c}):c()}}),St={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},It=z("path",{d:"M336 64h32a48 48 0 0 1 48 48v320a48 48 0 0 1-48 48H144a48 48 0 0 1-48-48V112a48 48 0 0 1 48-48h32",fill:"none",stroke:"currentColor","stroke-linejoin":"round","stroke-width":"32"},null,-1),Rt=z("rect",{x:"176",y:"32",width:"160",height:"64",rx:"26.13",ry:"26.13",fill:"none",stroke:"currentColor","stroke-linejoin":"round","stroke-width":"32"},null,-1),$t=[It,Rt],Pt=S({name:"ClipboardOutline",render:function(o,t){return j(),E("svg",St,$t)}}),Ot={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Tt=z("rect",{x:"48",y:"48",width:"176",height:"176",rx:"20",ry:"20",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Mt=z("rect",{x:"288",y:"48",width:"176",height:"176",rx:"20",ry:"20",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Nt=z("rect",{x:"48",y:"288",width:"176",height:"176",rx:"20",ry:"20",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Bt=z("rect",{x:"288",y:"288",width:"176",height:"176",rx:"20",ry:"20",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),jt=[Tt,Mt,Nt,Bt],At=S({name:"GridOutline",render:function(o,t){return j(),E("svg",Ot,jt)}}),Ht={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Lt=He('<rect x="80" y="80" width="352" height="352" rx="48" ry="48" fill="none" stroke="currentColor" stroke-linejoin="round" stroke-width="32"></rect><rect x="144" y="144" width="224" height="224" rx="16" ry="16" fill="none" stroke="currentColor" stroke-linejoin="round" stroke-width="32"></rect><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M256 80V48"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M336 80V48"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M176 80V48"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M256 464v-32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M336 464v-32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M176 464v-32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M432 256h32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M432 336h32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M432 176h32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M48 256h32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M48 336h32"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="32" d="M48 176h32"></path>',14),Et=[Lt],Ft=S({name:"HardwareChipOutline",render:function(o,t){return j(),E("svg",Ht,Et)}}),Kt={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Vt=z("path",{d:"M80 212v236a16 16 0 0 0 16 16h96V328a24 24 0 0 1 24-24h80a24 24 0 0 1 24 24v136h96a16 16 0 0 0 16-16V212",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Dt=z("path",{d:"M480 256L266.89 52c-5-5.28-16.69-5.34-21.78 0L32 256",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Ut=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32",d:"M400 179V64h-48v69"},null,-1),Gt=[Vt,Dt,Ut],qt=S({name:"HomeOutline",render:function(o,t){return j(),E("svg",Kt,Gt)}}),Wt={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Yt=z("path",{d:"M218.1 167.17c0 13 0 25.6 4.1 37.4c-43.1 50.6-156.9 184.3-167.5 194.5a20.17 20.17 0 0 0-6.7 15c0 8.5 5.2 16.7 9.6 21.3c6.6 6.9 34.8 33 40 28c15.4-15 18.5-19 24.8-25.2c9.5-9.3-1-28.3 2.3-36s6.8-9.2 12.5-10.4s15.8 2.9 23.7 3c8.3.1 12.8-3.4 19-9.2c5-4.6 8.6-8.9 8.7-15.6c.2-9-12.8-20.9-3.1-30.4s23.7 6.2 34 5s22.8-15.5 24.1-21.6s-11.7-21.8-9.7-30.7c.7-3 6.8-10 11.4-11s25 6.9 29.6 5.9c5.6-1.2 12.1-7.1 17.4-10.4c15.5 6.7 29.6 9.4 47.7 9.4c68.5 0 124-53.4 124-119.2S408.5 48 340 48s-121.9 53.37-121.9 119.17zM400 144a32 32 0 1 1-32-32a32 32 0 0 1 32 32z",fill:"none",stroke:"currentColor","stroke-linejoin":"round","stroke-width":"32"},null,-1),Xt=[Yt],Zt=S({name:"KeyOutline",render:function(o,t){return j(),E("svg",Wt,Xt)}}),Jt={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Qt=z("path",{d:"M304 336v40a40 40 0 0 1-40 40H104a40 40 0 0 1-40-40V136a40 40 0 0 1 40-40h152c22.09 0 48 17.91 48 40v40",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),er=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32",d:"M368 336l80-80l-80-80"},null,-1),or=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32",d:"M176 256h256"},null,-1),tr=[Qt,er,or],rr=S({name:"LogOutOutline",render:function(o,t){return j(),E("svg",Jt,tr)}}),nr={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},ir=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-miterlimit":"10","stroke-width":"32",d:"M80 160h352"},null,-1),lr=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-miterlimit":"10","stroke-width":"32",d:"M80 256h352"},null,-1),ar=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-miterlimit":"10","stroke-width":"32",d:"M80 352h352"},null,-1),sr=[ir,lr,ar],cr=S({name:"MenuOutline",render:function(o,t){return j(),E("svg",nr,sr)}}),dr={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},ur=z("path",{d:"M160 136c0-30.62 4.51-61.61 16-88C99.57 81.27 48 159.32 48 248c0 119.29 96.71 216 216 216c88.68 0 166.73-51.57 200-128c-26.39 11.49-57.38 16-88 16c-119.29 0-216-96.71-216-216z",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),hr=[ur],vr=S({name:"MoonOutline",render:function(o,t){return j(),E("svg",dr,hr)}}),mr={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},pr=z("path",{d:"M427.68 351.43C402 320 383.87 304 383.87 217.35C383.87 138 343.35 109.73 310 96c-4.43-1.82-8.6-6-9.95-10.55C294.2 65.54 277.8 48 256 48s-38.21 17.55-44 37.47c-1.35 4.6-5.52 8.71-9.95 10.53c-33.39 13.75-73.87 41.92-73.87 121.35C128.13 304 110 320 84.32 351.43C73.68 364.45 83 384 101.61 384h308.88c18.51 0 27.77-19.61 17.19-32.57z",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),fr=z("path",{d:"M320 384v16a64 64 0 0 1-128 0v-16",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),gr=[pr,fr],br=S({name:"NotificationsOutline",render:function(o,t){return j(),E("svg",mr,gr)}}),xr={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Cr=z("path",{d:"M402 168c-2.93 40.67-33.1 72-66 72s-63.12-31.32-66-72c-3-42.31 26.37-72 66-72s69 30.46 66 72z",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),wr=z("path",{d:"M336 304c-65.17 0-127.84 32.37-143.54 95.41c-2.08 8.34 3.15 16.59 11.72 16.59h263.65c8.57 0 13.77-8.25 11.72-16.59C463.85 335.36 401.18 304 336 304z",fill:"none",stroke:"currentColor","stroke-miterlimit":"10","stroke-width":"32"},null,-1),kr=z("path",{d:"M200 185.94c-2.34 32.48-26.72 58.06-53 58.06s-50.7-25.57-53-58.06C91.61 152.15 115.34 128 147 128s55.39 24.77 53 57.94z",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),yr=z("path",{d:"M206 306c-18.05-8.27-37.93-11.45-59-11.45c-52 0-102.1 25.85-114.65 76.2c-1.65 6.66 2.53 13.25 9.37 13.25H154",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-miterlimit":"10","stroke-width":"32"},null,-1),_r=[Cr,wr,kr,yr],zr=S({name:"PeopleOutline",render:function(o,t){return j(),E("svg",xr,_r)}}),Sr={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Ir=z("rect",{x:"96",y:"48",width:"320",height:"416",rx:"48",ry:"48",fill:"none",stroke:"currentColor","stroke-linejoin":"round","stroke-width":"32"},null,-1),Rr=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32",d:"M176 128h160"},null,-1),$r=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32",d:"M176 208h160"},null,-1),Pr=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32",d:"M176 288h80"},null,-1),Or=[Ir,Rr,$r,Pr],Tr=S({name:"ReaderOutline",render:function(o,t){return j(),E("svg",Sr,Or)}}),Mr={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Nr=z("path",{d:"M262.29 192.31a64 64 0 1 0 57.4 57.4a64.13 64.13 0 0 0-57.4-57.4zM416.39 256a154.34 154.34 0 0 1-1.53 20.79l45.21 35.46a10.81 10.81 0 0 1 2.45 13.75l-42.77 74a10.81 10.81 0 0 1-13.14 4.59l-44.9-18.08a16.11 16.11 0 0 0-15.17 1.75A164.48 164.48 0 0 1 325 400.8a15.94 15.94 0 0 0-8.82 12.14l-6.73 47.89a11.08 11.08 0 0 1-10.68 9.17h-85.54a11.11 11.11 0 0 1-10.69-8.87l-6.72-47.82a16.07 16.07 0 0 0-9-12.22a155.3 155.3 0 0 1-21.46-12.57a16 16 0 0 0-15.11-1.71l-44.89 18.07a10.81 10.81 0 0 1-13.14-4.58l-42.77-74a10.8 10.8 0 0 1 2.45-13.75l38.21-30a16.05 16.05 0 0 0 6-14.08c-.36-4.17-.58-8.33-.58-12.5s.21-8.27.58-12.35a16 16 0 0 0-6.07-13.94l-38.19-30A10.81 10.81 0 0 1 49.48 186l42.77-74a10.81 10.81 0 0 1 13.14-4.59l44.9 18.08a16.11 16.11 0 0 0 15.17-1.75A164.48 164.48 0 0 1 187 111.2a15.94 15.94 0 0 0 8.82-12.14l6.73-47.89A11.08 11.08 0 0 1 213.23 42h85.54a11.11 11.11 0 0 1 10.69 8.87l6.72 47.82a16.07 16.07 0 0 0 9 12.22a155.3 155.3 0 0 1 21.46 12.57a16 16 0 0 0 15.11 1.71l44.89-18.07a10.81 10.81 0 0 1 13.14 4.58l42.77 74a10.8 10.8 0 0 1-2.45 13.75l-38.21 30a16.05 16.05 0 0 0-6.05 14.08c.33 4.14.55 8.3.55 12.47z",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Br=[Nr],jr=S({name:"SettingsOutline",render:function(o,t){return j(),E("svg",Mr,Br)}}),Ar={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Hr=z("path",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32",d:"M336 176L225.2 304L176 255.8"},null,-1),Lr=z("path",{d:"M463.1 112.37C373.68 96.33 336.71 84.45 256 48c-80.71 36.45-117.68 48.33-207.1 64.37C32.7 369.13 240.58 457.79 256 464c15.42-6.21 223.3-94.87 207.1-351.63z",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Er=[Hr,Lr],Fr=S({name:"ShieldCheckmarkOutline",render:function(o,t){return j(),E("svg",Ar,Er)}}),Kr={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Vr=z("rect",{x:"64",y:"320",width:"48",height:"160",rx:"8",ry:"8",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Dr=z("rect",{x:"288",y:"224",width:"48",height:"256",rx:"8",ry:"8",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Ur=z("rect",{x:"400",y:"112",width:"48",height:"368",rx:"8",ry:"8",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),Gr=z("rect",{x:"176",y:"32",width:"48",height:"448",rx:"8",ry:"8",fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-linejoin":"round","stroke-width":"32"},null,-1),qr=[Vr,Dr,Ur,Gr],Wr=S({name:"StatsChartOutline",render:function(o,t){return j(),E("svg",Kr,qr)}}),Yr={xmlns:"http://www.w3.org/2000/svg","xmlns:xlink":"http://www.w3.org/1999/xlink",viewBox:"0 0 512 512"},Xr=He('<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M256 48v48"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M256 416v48"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M403.08 108.92l-33.94 33.94"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M142.86 369.14l-33.94 33.94"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M464 256h-48"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M96 256H48"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M403.08 403.08l-33.94-33.94"></path><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32" d="M142.86 142.86l-33.94-33.94"></path><circle cx="256" cy="256" r="80" fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" stroke-width="32"></circle>',9),Zr=[Xr],Jr=S({name:"SunnyOutline",render:function(o,t){return j(),E("svg",Yr,Zr)}}),Qr=S({__name:"SideMenu",setup(e){const o=ye(),t=Le(),i=_o(),s=ke(),n={"grid-outline":At,"settings-outline":jr,"people-outline":zr,"hardware-chip-outline":Ft,"document-text-outline":Vo,"notifications-outline":br,"reader-outline":Tr,"warning-outline":Ko,"shield-checkmark-outline":Fr,"stats-chart-outline":Wr,"clipboard-outline":Pt};function u(y){const C=n[y];if(C)return()=>d(Z,null,{default:()=>d(C)})}function c(y){return y.map(C=>({label:C.label,key:C.key,icon:C.icon?u(C.icon):void 0,children:C.children?c(C.children):void 0}))}const l=x(()=>c(i.visibleMenus)),h=F(null);function M(){const y=t.path,C=P(i.visibleMenus);let f=null;for(const I of C)y.startsWith(I)&&(!f||I.length>f.length)&&(f=I);return f||"/"}function P(y){const C=[];for(const f of y)C.push(f.key),f.children&&C.push(...P(f.children));return C}h.value=M();function p(y){o.push(y)}return(y,C)=>(j(),J(B(zt),{value:h.value,collapsed:B(s).sidebarCollapsed,"collapsed-width":64,"collapsed-icon-size":22,options:l.value,indent:24,"onUpdate:value":p},null,8,["value","collapsed","options"]))}}),en=he(Qr,[["__scopeId","data-v-765a7b2f"]]),on={class:"breadcrumb-wrapper"},tn=S({__name:"Breadcrumb",setup(e){const o=Le(),t=ye(),i=x(()=>o.matched.filter(u=>{var c;return(c=u.meta)==null?void 0:c.title}).map(u=>{var c;return{title:(c=u.meta)==null?void 0:c.title,path:u.path}}));function s(n){n!==o.path&&t.push(n)}return(n,u)=>(j(),E("div",on,[K(B(ot),{separator:">"},{default:D(()=>[K(B(Te),null,{default:D(()=>[K(B(Z),{component:B(qt),size:16},null,8,["component"])]),_:1}),(j(!0),E(Ae,null,zo(i.value,(c,l)=>(j(),J(B(Te),{key:l,onClick:h=>s(c.path)},{default:D(()=>[So(Io(c.title),1)]),_:2},1032,["onClick"]))),128))]),_:1})]))}}),rn=he(tn,[["__scopeId","data-v-b31f7704"]]),nn={class:"top-header"},ln={class:"top-header-left"},an={class:"top-header-right"},sn=S({__name:"TopHeader",setup(e){const o=ye(),t=ke(),i=Ro(),s=[{label:"修改密码",key:"change-password",icon:()=>d(Z,null,{default:()=>d(Zt)})},{type:"divider",key:"d1"},{label:"退出登录",key:"logout",icon:()=>d(Z,null,{default:()=>d(rr)})}];function n(c){c==="logout"?(i.logout(),o.push("/login")):c==="change-password"&&o.push("/change-password")}function u(){var c;return d("div",{style:{display:"flex",alignItems:"center",gap:"8px"}},[d(Jo,{size:"small",round:!0,style:{backgroundColor:"var(--primary-color)"}},{default:()=>{var l;return(((l=i.user)==null?void 0:l.displayName)||i.username||"U").charAt(0).toUpperCase()}}),d("span",null,((c=i.user)==null?void 0:c.displayName)||i.username)])}return(c,l)=>(j(),E("div",nn,[z("div",ln,[K(B(fe),{text:"",style:{"font-size":"20px"},onClick:l[0]||(l[0]=h=>B(t).toggleSidebar())},{default:D(()=>[K(B(Z),{component:B(cr)},null,8,["component"])]),_:1}),K(rn)]),z("div",an,[K(B(Do),{align:"center"},{default:D(()=>[K(B(fe),{text:"",style:{"font-size":"20px"},onClick:l[1]||(l[1]=h=>B(t).toggleTheme())},{default:D(()=>[B(t).theme==="light"?(j(),J(B(Z),{key:0,component:B(vr)},null,8,["component"])):(j(),J(B(Z),{key:1,component:B(Jr)},null,8,["component"]))]),_:1}),K(B(Ue),{options:s,onSelect:n,trigger:"click"},{default:D(()=>[K(B(fe),{text:""},{default:D(()=>[(j(),J($o(u())))]),_:1})]),_:1})]),_:1})])]))}}),cn=he(sn,[["__scopeId","data-v-28756eac"]]),dn={class:"sidebar-logo"},un={key:0,class:"logo-text"},hn={key:1,class:"logo-text-collapsed"},vn=S({__name:"AppLayout",setup(e){const o=ke();return(t,i)=>(j(),J(B(Me),{"has-sider":"",position:"absolute",style:{top:"0",left:"0",right:"0",bottom:"0"}},{default:D(()=>[K(B(pt),{bordered:"",collapsed:B(o).sidebarCollapsed,"collapsed-width":64,width:240,"native-scrollbar":!1,"collapse-mode":"width"},{default:D(()=>[z("div",dn,[B(o).sidebarCollapsed?(j(),E("span",hn,"GL")):(j(),E("span",un,"GoLog"))]),K(en)]),_:1},8,["collapsed"]),K(B(Me),null,{default:D(()=>[K(B(dt),{bordered:""},{default:D(()=>[K(cn)]),_:1}),K(B(at),{"native-scrollbar":!1,"content-style":"padding: 16px 24px; min-height: calc(100vh - 56px)"},{default:D(()=>[Po(t.$slots,"default",{},void 0,!0)]),_:3})]),_:3})]),_:3}))}}),mn=he(vn,[["__scopeId","data-v-103f2ac4"]]),Tn=S({__name:"DefaultLayout",setup(e){return(o,t)=>{const i=Oo("router-view");return j(),J(mn,null,{default:D(()=>[K(i)]),_:1})}}});export{Tn as default};
