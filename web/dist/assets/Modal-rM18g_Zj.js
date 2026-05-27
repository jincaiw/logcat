import{b3 as G,b4 as R,aQ as ie,aR as Q,aT as Xe,bA as U,W as Ye,s as _,t as T,v as S,w as N,as as We,h as Ke,a3 as se,al as d,by as J,a6 as qe,N as F,X as Ue,aU as Ve,E as Ge,aJ as Je,bD as ee,e as Be,bz as ce,aq as te,aM as de,bk as X,aN as ue,a$ as V,aK as Qe}from"./index-qenDAa9H.js";import{a as Ze,N as et,c as tt}from"./Card-Dq-y6vVg.js";import{w as I,v as L,q as ot,u as Se,C as Z,a as nt,A as fe,N as it,F as Oe,J as st,L as Me,n as lt,M as at,S as rt,k as ct,r as dt,e as W,o as ut}from"./Scrollbar-DUkKw_1i.js";import{r as j}from"./create-CSfaedQP.js";import{B as ve}from"./Button-Bkv4AwgB.js";import{E as ft,W as vt,S as ht,I as he}from"./Warning-C_UQMDDl.js";import{g as gt,F as mt,b as pt,k as oe,j as $e,m as bt,e as Ct,n as yt,f as wt,L as kt,z as xt,l as Rt}from"./Popover-BJ4Ic_2j.js";import{i as Fe,h as Te}from"./utils-DBY6nFc7.js";import{e as Pt}from"./Select-BSIolqJY.js";const H=R(null);function ge(e){if(e.clientX>0||e.clientY>0)H.value={x:e.clientX,y:e.clientY};else{const{target:t}=e;if(t instanceof Element){const{left:o,top:a,width:u,height:f}=t.getBoundingClientRect();o>0||a>0?H.value={x:o+u/2,y:a+f/2}:H.value={x:0,y:0}}else H.value=null}}let K=0,me=!0;function Bt(){if(!Fe)return G(R(null));K===0&&I("click",document,ge,!0);const e=()=>{K+=1};return me&&(me=Te())?(ie(e),Q(()=>{K-=1,K===0&&L("click",document,ge,!0)})):e(),G(H)}const St=R(void 0);let q=0;function pe(){St.value=Date.now()}let be=!0;function Ot(e){if(!Fe)return G(R(!1));const t=R(!1);let o=null;function a(){o!==null&&window.clearTimeout(o)}function u(){a(),t.value=!0,o=window.setTimeout(()=>{t.value=!1},e)}q===0&&I("click",window,pe,!0);const f=()=>{q+=1,I("click",window,u,!0)};return be&&(be=Te())?(ie(f),Q(()=>{q-=1,q===0&&L("click",window,pe,!0),L("click",window,u,!0),a()})):f(),G(t)}const le=R(!1);function Ce(){le.value=!0}function ye(){le.value=!1}let D=0;function Mt(){return ot&&(ie(()=>{D||(window.addEventListener("compositionstart",Ce),window.addEventListener("compositionend",ye)),D++}),Q(()=>{D<=1?(window.removeEventListener("compositionstart",Ce),window.removeEventListener("compositionend",ye),D=0):D--})),le}let A=0,we="",ke="",xe="",Re="";const Pe=R("0px");function $t(e){if(typeof document>"u")return;const t=document.documentElement;let o,a=!1;const u=()=>{t.style.marginRight=we,t.style.overflow=ke,t.style.overflowX=xe,t.style.overflowY=Re,Pe.value="0px"};Xe(()=>{o=U(e,f=>{if(f){if(!A){const v=window.innerWidth-t.offsetWidth;v>0&&(we=t.style.marginRight,t.style.marginRight=`${v}px`,Pe.value=`${v}px`),ke=t.style.overflow,xe=t.style.overflowX,Re=t.style.overflowY,t.style.overflow="hidden",t.style.overflowX="hidden",t.style.overflowY="hidden"}a=!0,A++}else A--,A||u(),a=!1},{immediate:!0})}),Q(()=>{o==null||o(),a&&(A--,A||u(),a=!1)})}const Ft=Ye("n-dialog-provider"),ae={icon:Function,type:{type:String,default:"default"},title:[String,Function],closable:{type:Boolean,default:!0},negativeText:String,positiveText:String,positiveButtonProps:Object,negativeButtonProps:Object,content:[String,Function],action:Function,showIcon:{type:Boolean,default:!0},loading:Boolean,bordered:Boolean,iconPlacement:String,titleClass:[String,Array],titleStyle:[String,Object],contentClass:[String,Array],contentStyle:[String,Object],actionClass:[String,Array],actionStyle:[String,Object],onPositiveClick:Function,onNegativeClick:Function,onClose:Function,closeFocusable:Boolean},Tt=Se(ae),Et=_([T("dialog",`
 --n-icon-margin: var(--n-icon-margin-top) var(--n-icon-margin-right) var(--n-icon-margin-bottom) var(--n-icon-margin-left);
 word-break: break-word;
 line-height: var(--n-line-height);
 position: relative;
 background: var(--n-color);
 color: var(--n-text-color);
 box-sizing: border-box;
 margin: auto;
 border-radius: var(--n-border-radius);
 padding: var(--n-padding);
 transition: 
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `,[S("icon",`
 color: var(--n-icon-color);
 `),N("bordered",`
 border: var(--n-border);
 `),N("icon-top",[S("close",`
 margin: var(--n-close-margin);
 `),S("icon",`
 margin: var(--n-icon-margin);
 `),S("content",`
 text-align: center;
 `),S("title",`
 justify-content: center;
 `),S("action",`
 justify-content: center;
 `)]),N("icon-left",[S("icon",`
 margin: var(--n-icon-margin);
 `),N("closable",[S("title",`
 padding-right: calc(var(--n-close-size) + 6px);
 `)])]),S("close",`
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 z-index: 1;
 `),S("content",`
 font-size: var(--n-font-size);
 margin: var(--n-content-margin);
 position: relative;
 word-break: break-word;
 `,[N("last","margin-bottom: 0;")]),S("action",`
 display: flex;
 justify-content: flex-end;
 `,[_("> *:not(:last-child)",`
 margin-right: var(--n-action-space);
 `)]),S("icon",`
 font-size: var(--n-icon-size);
 transition: color .3s var(--n-bezier);
 `),S("title",`
 transition: color .3s var(--n-bezier);
 display: flex;
 align-items: center;
 font-size: var(--n-title-font-size);
 font-weight: var(--n-title-font-weight);
 color: var(--n-title-text-color);
 `),T("dialog-icon-container",`
 display: flex;
 justify-content: center;
 `)]),We(T("dialog",`
 width: 446px;
 max-width: calc(100vw - 32px);
 `)),T("dialog",[Ke(`
 width: 446px;
 max-width: calc(100vw - 32px);
 `)])]),zt={default:()=>d(he,null),info:()=>d(he,null),success:()=>d(ht,null),warning:()=>d(vt,null),error:()=>d(ft,null)},jt=se({name:"Dialog",alias:["NimbusConfirmCard","Confirm"],props:Object.assign(Object.assign({},J.props),ae),slots:Object,setup(e){const{mergedComponentPropsRef:t,mergedClsPrefixRef:o,inlineThemeDisabled:a,mergedRtlRef:u}=Oe(e),f=st("Dialog",u,o),v=F(()=>{var h,g;const{iconPlacement:k}=e;return k||((g=(h=t==null?void 0:t.value)===null||h===void 0?void 0:h.Dialog)===null||g===void 0?void 0:g.iconPlacement)||"left"});function b(h){const{onPositiveClick:g}=e;g&&g(h)}function l(h){const{onNegativeClick:g}=e;g&&g(h)}function x(){const{onClose:h}=e;h&&h()}const w=J("Dialog","-dialog",Et,qe,e,o),C=F(()=>{const{type:h}=e,g=v.value,{common:{cubicBezierEaseInOut:k},self:{fontSize:P,lineHeight:y,border:r,titleTextColor:O,textColor:M,color:p,closeBorderRadius:n,closeColorHover:s,closeColorPressed:i,closeIconColor:m,closeIconColorHover:B,closeIconColorPressed:$,closeIconSize:E,borderRadius:z,titleFontWeight:Ee,titleFontSize:ze,padding:je,iconSize:Ae,actionSpace:Ne,contentMargin:Ie,closeSize:Le,[g==="top"?"iconMarginIconTop":"iconMargin"]:De,[g==="top"?"closeMarginIconTop":"closeMargin"]:He,[Ue("iconColor",h)]:_e}}=w.value,Y=lt(De);return{"--n-font-size":P,"--n-icon-color":_e,"--n-bezier":k,"--n-close-margin":He,"--n-icon-margin-top":Y.top,"--n-icon-margin-right":Y.right,"--n-icon-margin-bottom":Y.bottom,"--n-icon-margin-left":Y.left,"--n-icon-size":Ae,"--n-close-size":Le,"--n-close-icon-size":E,"--n-close-border-radius":n,"--n-close-color-hover":s,"--n-close-color-pressed":i,"--n-close-icon-color":m,"--n-close-icon-color-hover":B,"--n-close-icon-color-pressed":$,"--n-color":p,"--n-text-color":M,"--n-border-radius":z,"--n-padding":je,"--n-line-height":y,"--n-border":r,"--n-content-margin":Ie,"--n-title-font-size":ze,"--n-title-font-weight":Ee,"--n-title-text-color":O,"--n-action-space":Ne}}),c=a?Me("dialog",F(()=>`${e.type[0]}${v.value[0]}`),C,e):void 0;return{mergedClsPrefix:o,rtlEnabled:f,mergedIconPlacement:v,mergedTheme:w,handlePositiveClick:b,handleNegativeClick:l,handleCloseClick:x,cssVars:a?void 0:C,themeClass:c==null?void 0:c.themeClass,onRender:c==null?void 0:c.onRender}},render(){var e;const{bordered:t,mergedIconPlacement:o,cssVars:a,closable:u,showIcon:f,title:v,content:b,action:l,negativeText:x,positiveText:w,positiveButtonProps:C,negativeButtonProps:c,handlePositiveClick:h,handleNegativeClick:g,mergedTheme:k,loading:P,type:y,mergedClsPrefix:r}=this;(e=this.onRender)===null||e===void 0||e.call(this);const O=f?d(nt,{clsPrefix:r,class:`${r}-dialog__icon`},{default:()=>Z(this.$slots.icon,p=>p||(this.icon?j(this.icon):zt[this.type]()))}):null,M=Z(this.$slots.action,p=>p||w||x||l?d("div",{class:[`${r}-dialog__action`,this.actionClass],style:this.actionStyle},p||(l?[j(l)]:[this.negativeText&&d(ve,Object.assign({theme:k.peers.Button,themeOverrides:k.peerOverrides.Button,ghost:!0,size:"small",onClick:g},c),{default:()=>j(this.negativeText)}),this.positiveText&&d(ve,Object.assign({theme:k.peers.Button,themeOverrides:k.peerOverrides.Button,size:"small",type:y==="default"?"primary":y,disabled:P,loading:P,onClick:h},C),{default:()=>j(this.positiveText)})])):null);return d("div",{class:[`${r}-dialog`,this.themeClass,this.closable&&`${r}-dialog--closable`,`${r}-dialog--icon-${o}`,t&&`${r}-dialog--bordered`,this.rtlEnabled&&`${r}-dialog--rtl`],style:a,role:"dialog"},u?Z(this.$slots.close,p=>{const n=[`${r}-dialog__close`,this.rtlEnabled&&`${r}-dialog--rtl`];return p?d("div",{class:n},p):d(it,{focusable:this.closeFocusable,clsPrefix:r,class:n,onClick:this.handleCloseClick})}):null,f&&o==="top"?d("div",{class:`${r}-dialog-icon-container`},O):null,d("div",{class:[`${r}-dialog__title`,this.titleClass],style:this.titleStyle},f&&o==="left"?O:null,fe(this.$slots.header,()=>[j(v)])),d("div",{class:[`${r}-dialog__content`,M?"":`${r}-dialog__content--last`,this.contentClass],style:this.contentStyle},fe(this.$slots.default,()=>[j(b)])),M)}}),ne="n-draggable";function At(e,t){let o;const a=F(()=>e.value!==!1),u=F(()=>a.value?ne:""),f=F(()=>{const l=e.value;return l===!0||l===!1?!0:l?l.bounds!=="none":!0});function v(l){const x=l.querySelector(`.${ne}`);if(!x||!u.value)return;let w=0,C=0,c=0,h=0,g=0,k=0,P,y=null,r=null;function O(s){s.preventDefault(),P=s;const{x:i,y:m,right:B,bottom:$}=l.getBoundingClientRect();C=i,h=m,w=window.innerWidth-B,c=window.innerHeight-$;const{left:E,top:z}=l.style;g=+z.slice(0,-2),k=+E.slice(0,-2)}function M(){r&&(l.style.top=`${r.y}px`,l.style.left=`${r.x}px`,r=null),y=null}function p(s){if(!P)return;const{clientX:i,clientY:m}=P;let B=s.clientX-i,$=s.clientY-m;f.value&&(B>w?B=w:-B>C&&(B=-C),$>c?$=c:-$>h&&($=-h));const E=B+k,z=$+g;r={x:E,y:z},y||(y=requestAnimationFrame(M))}function n(){P=void 0,y&&(cancelAnimationFrame(y),y=null),r&&(l.style.top=`${r.y}px`,l.style.left=`${r.x}px`,r=null),t.onEnd(l)}I("mousedown",x,O),I("mousemove",window,p),I("mouseup",window,n),o=()=>{y&&cancelAnimationFrame(y),L("mousedown",x,O),L("mousemove",window,p),L("mouseup",window,n)}}function b(){o&&(o(),o=void 0)}return Ve(b),{stopDrag:b,startDrag:v,draggableRef:a,draggableClassRef:u}}const re=Object.assign(Object.assign({},Ze),ae),Nt=Se(re),It=se({name:"ModalBody",inheritAttrs:!1,slots:Object,props:Object.assign(Object.assign({show:{type:Boolean,required:!0},preset:String,displayDirective:{type:String,required:!0},trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},blockScroll:Boolean,draggable:{type:[Boolean,Object],default:!1},maskHidden:Boolean},re),{renderMask:Function,onClickoutside:Function,onBeforeLeave:{type:Function,required:!0},onAfterLeave:{type:Function,required:!0},onPositiveClick:{type:Function,required:!0},onNegativeClick:{type:Function,required:!0},onClose:{type:Function,required:!0},onAfterEnter:Function,onEsc:Function}),setup(e){const t=R(null),o=R(null),a=R(e.show),u=R(null),f=R(null),v=te($e);let b=null;U(X(e,"show"),i=>{i&&(b=v.getMousePosition())},{immediate:!0});const{stopDrag:l,startDrag:x,draggableRef:w,draggableClassRef:C}=At(X(e,"draggable"),{onEnd:i=>{k(i)}}),c=F(()=>ue([e.titleClass,C.value])),h=F(()=>ue([e.headerClass,C.value]));U(X(e,"show"),i=>{i&&(a.value=!0)}),$t(F(()=>e.blockScroll&&a.value));function g(){if(v.transformOriginRef.value==="center")return"";const{value:i}=u,{value:m}=f;if(i===null||m===null)return"";if(o.value){const B=o.value.containerScrollTop;return`${i}px ${m+B}px`}return""}function k(i){if(v.transformOriginRef.value==="center"||!b||!o.value)return;const m=o.value.containerScrollTop,{offsetLeft:B,offsetTop:$}=i,E=b.y,z=b.x;u.value=-(B-z),f.value=-($-E-m),i.style.transformOrigin=g()}function P(i){de(()=>{k(i)})}function y(i){i.style.transformOrigin=g(),e.onBeforeLeave()}function r(i){const m=i;w.value&&x(m),e.onAfterEnter&&e.onAfterEnter(m)}function O(){a.value=!1,u.value=null,f.value=null,l(),e.onAfterLeave()}function M(){const{onClose:i}=e;i&&i()}function p(){e.onNegativeClick()}function n(){e.onPositiveClick()}const s=R(null);return U(s,i=>{i&&de(()=>{const m=i.el;m&&t.value!==m&&(t.value=m)})}),V(bt,t),V(Ct,null),V(yt,null),{mergedTheme:v.mergedThemeRef,appear:v.appearRef,isMounted:v.isMountedRef,mergedClsPrefix:v.mergedClsPrefixRef,bodyRef:t,scrollbarRef:o,draggableClass:C,displayed:a,childNodeRef:s,cardHeaderClass:h,dialogTitleClass:c,handlePositiveClick:n,handleNegativeClick:p,handleCloseClick:M,handleAfterEnter:r,handleAfterLeave:O,handleBeforeLeave:y,handleEnter:P}},render(){const{$slots:e,$attrs:t,handleEnter:o,handleAfterEnter:a,handleAfterLeave:u,handleBeforeLeave:f,preset:v,mergedClsPrefix:b}=this;let l=null;if(!v){if(l=gt("default",e.default,{draggableClass:this.draggableClass}),!l){at("modal","default slot is empty");return}l=Ge(l),l.props=Je({class:`${b}-modal`},t,l.props||{})}return this.displayDirective==="show"||this.displayed||this.show?ee(d("div",{role:"none",class:[`${b}-modal-body-wrapper`,this.maskHidden&&`${b}-modal-body-wrapper--mask-hidden`]},d(rt,{ref:"scrollbarRef",theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,contentClass:`${b}-modal-scroll-content`},{default:()=>{var x;return[(x=this.renderMask)===null||x===void 0?void 0:x.call(this),d(mt,{disabled:!this.trapFocus||this.maskHidden,active:this.show,onEsc:this.onEsc,autoFocus:this.autoFocus},{default:()=>{var w;return d(Be,{name:"fade-in-scale-up-transition",appear:(w=this.appear)!==null&&w!==void 0?w:this.isMounted,onEnter:o,onAfterEnter:a,onAfterLeave:u,onBeforeLeave:f},{default:()=>{const C=[[ce,this.show]],{onClickoutside:c}=this;return c&&C.push([pt,this.onClickoutside,void 0,{capture:!0}]),ee(this.preset==="confirm"||this.preset==="dialog"?d(jt,Object.assign({},this.$attrs,{class:[`${b}-modal`,this.$attrs.class],ref:"bodyRef",theme:this.mergedTheme.peers.Dialog,themeOverrides:this.mergedTheme.peerOverrides.Dialog},oe(this.$props,Tt),{titleClass:this.dialogTitleClass,"aria-modal":"true"}),e):this.preset==="card"?d(et,Object.assign({},this.$attrs,{ref:"bodyRef",class:[`${b}-modal`,this.$attrs.class],theme:this.mergedTheme.peers.Card,themeOverrides:this.mergedTheme.peerOverrides.Card},oe(this.$props,tt),{headerClass:this.cardHeaderClass,"aria-modal":"true",role:"dialog"}),e):this.childNodeRef=l,C)}})}})]}})),[[ce,this.displayDirective==="if"||this.displayed||this.show]]):null}}),Lt=_([T("modal-container",`
 position: fixed;
 left: 0;
 top: 0;
 height: 0;
 width: 0;
 display: flex;
 `),T("modal-mask",`
 position: fixed;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 background-color: rgba(0, 0, 0, .4);
 `,[ct({enterDuration:".25s",leaveDuration:".25s",enterCubicBezier:"var(--n-bezier-ease-out)",leaveCubicBezier:"var(--n-bezier-ease-out)"})]),T("modal-body-wrapper",`
 position: fixed;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 overflow: visible;
 `,[T("modal-scroll-content",`
 min-height: 100%;
 display: flex;
 position: relative;
 `),N("mask-hidden","pointer-events: none;",[T("modal-scroll-content",[_("> *",`
 pointer-events: all;
 `)])])]),T("modal",`
 position: relative;
 align-self: center;
 color: var(--n-text-color);
 margin: auto;
 box-shadow: var(--n-box-shadow);
 `,[wt({duration:".25s",enterScale:".5"}),_(`.${ne}`,`
 cursor: move;
 user-select: none;
 `)])]),Dt=Object.assign(Object.assign(Object.assign(Object.assign({},J.props),{show:Boolean,showMask:{type:Boolean,default:!0},maskClosable:{type:Boolean,default:!0},preset:String,to:[String,Object],displayDirective:{type:String,default:"if"},transformOrigin:{type:String,default:"mouse"},zIndex:Number,autoFocus:{type:Boolean,default:!0},trapFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0}}),re),{draggable:[Boolean,Object],onEsc:Function,"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onBeforeLeave:Function,onAfterLeave:Function,onClose:Function,onPositiveClick:Function,onNegativeClick:Function,onMaskClick:Function,internalDialog:Boolean,internalModal:Boolean,internalAppear:{type:Boolean,default:void 0},overlayStyle:[String,Object],onBeforeHide:Function,onAfterHide:Function,onHide:Function,unstableShowMask:{type:Boolean,default:void 0}}),Gt=se({name:"Modal",inheritAttrs:!1,props:Dt,slots:Object,setup(e){const t=R(null),{mergedClsPrefixRef:o,namespaceRef:a,inlineThemeDisabled:u}=Oe(e),f=J("Modal","-modal",Lt,Qe,e,o),v=Ot(64),b=Bt(),l=dt(),x=e.internalDialog?te(Ft,null):null,w=e.internalModal?te(Rt,null):null,C=Mt();function c(n){const{onUpdateShow:s,"onUpdate:show":i,onHide:m}=e;s&&W(s,n),i&&W(i,n),m&&!n&&m(n)}function h(){const{onClose:n}=e;n?Promise.resolve(n()).then(s=>{s!==!1&&c(!1)}):c(!1)}function g(){const{onPositiveClick:n}=e;n?Promise.resolve(n()).then(s=>{s!==!1&&c(!1)}):c(!1)}function k(){const{onNegativeClick:n}=e;n?Promise.resolve(n()).then(s=>{s!==!1&&c(!1)}):c(!1)}function P(){const{onBeforeLeave:n,onBeforeHide:s}=e;n&&W(n),s&&s()}function y(){const{onAfterLeave:n,onAfterHide:s}=e;n&&W(n),s&&s()}function r(n){var s;const{onMaskClick:i}=e;i&&i(n),e.maskClosable&&!((s=t.value)===null||s===void 0)&&s.contains(ut(n))&&c(!1)}function O(n){var s;(s=e.onEsc)===null||s===void 0||s.call(e),e.show&&e.closeOnEsc&&Pt(n)&&(C.value||c(!1))}V($e,{getMousePosition:()=>{const n=x||w;if(n){const{clickedRef:s,clickedPositionRef:i}=n;if(s.value&&i.value)return i.value}return v.value?b.value:null},mergedClsPrefixRef:o,mergedThemeRef:f,isMountedRef:l,appearRef:X(e,"internalAppear"),transformOriginRef:X(e,"transformOrigin")});const M=F(()=>{const{common:{cubicBezierEaseOut:n},self:{boxShadow:s,color:i,textColor:m}}=f.value;return{"--n-bezier-ease-out":n,"--n-box-shadow":s,"--n-color":i,"--n-text-color":m}}),p=u?Me("theme-class",void 0,M,e):void 0;return{mergedClsPrefix:o,namespace:a,isMounted:l,containerRef:t,presetProps:F(()=>oe(e,Nt)),handleEsc:O,handleAfterLeave:y,handleClickoutside:r,handleBeforeLeave:P,doUpdateShow:c,handleNegativeClick:k,handlePositiveClick:g,handleCloseClick:h,cssVars:u?void 0:M,themeClass:p==null?void 0:p.themeClass,onRender:p==null?void 0:p.onRender}},render(){const{mergedClsPrefix:e}=this;return d(kt,{to:this.to,show:this.show},{default:()=>{var t;(t=this.onRender)===null||t===void 0||t.call(this);const{showMask:o}=this;return ee(d("div",{role:"none",ref:"containerRef",class:[`${e}-modal-container`,this.themeClass,this.namespace],style:this.cssVars},d(It,Object.assign({style:this.overlayStyle},this.$attrs,{ref:"bodyWrapper",displayDirective:this.displayDirective,show:this.show,preset:this.preset,autoFocus:this.autoFocus,trapFocus:this.trapFocus,draggable:this.draggable,blockScroll:this.blockScroll,maskHidden:!o},this.presetProps,{onEsc:this.handleEsc,onClose:this.handleCloseClick,onNegativeClick:this.handleNegativeClick,onPositiveClick:this.handlePositiveClick,onBeforeLeave:this.handleBeforeLeave,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave,onClickoutside:o?void 0:this.handleClickoutside,renderMask:o?()=>{var a;return d(Be,{name:"fade-in-transition",key:"mask",appear:(a=this.internalAppear)!==null&&a!==void 0?a:this.isMounted},{default:()=>this.show?d("div",{"aria-hidden":!0,ref:"containerRef",class:`${e}-modal-mask`,onClick:this.handleClickoutside}):null})}:void 0}),this.$slots)),[[xt,{zIndex:this.zIndex,enabled:this.show}]])}})}});export{Gt as N};
