import{F,J as O,L as M,N as V,A as J,a as D,C as K,n as X}from"./Scrollbar-DUkKw_1i.js";import{f as q}from"./fade-in-height-expand.cssr-Dro2CEvB.js";import{a5 as G,J as Q,L as v,z as f,t as I,v as i,w as T,s as U,a3 as Y,al as t,by as $,N as R,b4 as Z,aJ as oo,X as c}from"./index-qenDAa9H.js";import{N as eo}from"./Button-Bkv4AwgB.js";import{E as ro,W as no,I as so,S as to}from"./Warning-C_UQMDDl.js";function lo(r){const{lineHeight:o,borderRadius:d,fontWeightStrong:b,baseColor:l,dividerColor:u,actionColor:y,textColor1:g,textColor2:s,closeColorHover:h,closeColorPressed:C,closeIconColor:m,closeIconColorHover:p,closeIconColorPressed:n,infoColor:e,successColor:x,warningColor:z,errorColor:S,fontSize:P}=r;return Object.assign(Object.assign({},Q),{fontSize:P,lineHeight:o,titleFontWeight:b,borderRadius:d,border:`1px solid ${u}`,color:y,titleTextColor:g,iconColor:s,contentTextColor:s,closeBorderRadius:d,closeColorHover:h,closeColorPressed:C,closeIconColor:m,closeIconColorHover:p,closeIconColorPressed:n,borderInfo:`1px solid ${v(l,f(e,{alpha:.25}))}`,colorInfo:v(l,f(e,{alpha:.08})),titleTextColorInfo:g,iconColorInfo:e,contentTextColorInfo:s,closeColorHoverInfo:h,closeColorPressedInfo:C,closeIconColorInfo:m,closeIconColorHoverInfo:p,closeIconColorPressedInfo:n,borderSuccess:`1px solid ${v(l,f(x,{alpha:.25}))}`,colorSuccess:v(l,f(x,{alpha:.08})),titleTextColorSuccess:g,iconColorSuccess:x,contentTextColorSuccess:s,closeColorHoverSuccess:h,closeColorPressedSuccess:C,closeIconColorSuccess:m,closeIconColorHoverSuccess:p,closeIconColorPressedSuccess:n,borderWarning:`1px solid ${v(l,f(z,{alpha:.33}))}`,colorWarning:v(l,f(z,{alpha:.08})),titleTextColorWarning:g,iconColorWarning:z,contentTextColorWarning:s,closeColorHoverWarning:h,closeColorPressedWarning:C,closeIconColorWarning:m,closeIconColorHoverWarning:p,closeIconColorPressedWarning:n,borderError:`1px solid ${v(l,f(S,{alpha:.25}))}`,colorError:v(l,f(S,{alpha:.08})),titleTextColorError:g,iconColorError:S,contentTextColorError:s,closeColorHoverError:h,closeColorPressedError:C,closeIconColorError:m,closeIconColorHoverError:p,closeIconColorPressedError:n})}const io={common:G,self:lo},ao=I("alert",`
 line-height: var(--n-line-height);
 border-radius: var(--n-border-radius);
 position: relative;
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 text-align: start;
 word-break: break-word;
`,[i("border",`
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 transition: border-color .3s var(--n-bezier);
 border: var(--n-border);
 pointer-events: none;
 `),T("closable",[I("alert-body",[i("title",`
 padding-right: 24px;
 `)])]),i("icon",{color:"var(--n-icon-color)"}),I("alert-body",{padding:"var(--n-padding)"},[i("title",{color:"var(--n-title-text-color)"}),i("content",{color:"var(--n-content-text-color)"})]),q({originalTransition:"transform .3s var(--n-bezier)",enterToProps:{transform:"scale(1)"},leaveToProps:{transform:"scale(0.9)"}}),i("icon",`
 position: absolute;
 left: 0;
 top: 0;
 align-items: center;
 justify-content: center;
 display: flex;
 width: var(--n-icon-size);
 height: var(--n-icon-size);
 font-size: var(--n-icon-size);
 margin: var(--n-icon-margin);
 `),i("close",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 position: absolute;
 right: 0;
 top: 0;
 margin: var(--n-close-margin);
 `),T("show-icon",[I("alert-body",{paddingLeft:"calc(var(--n-icon-margin-left) + var(--n-icon-size) + var(--n-icon-margin-right))"})]),T("right-adjust",[I("alert-body",{paddingRight:"calc(var(--n-close-size) + var(--n-padding) + 2px)"})]),I("alert-body",`
 border-radius: var(--n-border-radius);
 transition: border-color .3s var(--n-bezier);
 `,[i("title",`
 transition: color .3s var(--n-bezier);
 font-size: 16px;
 line-height: 19px;
 font-weight: var(--n-title-font-weight);
 `,[U("& +",[i("content",{marginTop:"9px"})])]),i("content",{transition:"color .3s var(--n-bezier)",fontSize:"var(--n-font-size)"})]),i("icon",{transition:"color .3s var(--n-bezier)"})]),co=Object.assign(Object.assign({},$.props),{title:String,showIcon:{type:Boolean,default:!0},type:{type:String,default:"default"},bordered:{type:Boolean,default:!0},closable:Boolean,onClose:Function,onAfterLeave:Function,onAfterHide:Function}),Co=Y({name:"Alert",inheritAttrs:!1,props:co,slots:Object,setup(r){const{mergedClsPrefixRef:o,mergedBorderedRef:d,inlineThemeDisabled:b,mergedRtlRef:l}=F(r),u=$("Alert","-alert",ao,io,r,o),y=O("Alert",l,o),g=R(()=>{const{common:{cubicBezierEaseInOut:n},self:e}=u.value,{fontSize:x,borderRadius:z,titleFontWeight:S,lineHeight:P,iconSize:H,iconMargin:E,iconMarginRtl:W,closeIconSize:w,closeBorderRadius:A,closeSize:B,closeMargin:_,closeMarginRtl:L,padding:j}=e,{type:a}=r,{left:k,right:N}=X(E);return{"--n-bezier":n,"--n-color":e[c("color",a)],"--n-close-icon-size":w,"--n-close-border-radius":A,"--n-close-color-hover":e[c("closeColorHover",a)],"--n-close-color-pressed":e[c("closeColorPressed",a)],"--n-close-icon-color":e[c("closeIconColor",a)],"--n-close-icon-color-hover":e[c("closeIconColorHover",a)],"--n-close-icon-color-pressed":e[c("closeIconColorPressed",a)],"--n-icon-color":e[c("iconColor",a)],"--n-border":e[c("border",a)],"--n-title-text-color":e[c("titleTextColor",a)],"--n-content-text-color":e[c("contentTextColor",a)],"--n-line-height":P,"--n-border-radius":z,"--n-font-size":x,"--n-title-font-weight":S,"--n-icon-size":H,"--n-icon-margin":E,"--n-icon-margin-rtl":W,"--n-close-size":B,"--n-close-margin":_,"--n-close-margin-rtl":L,"--n-padding":j,"--n-icon-margin-left":k,"--n-icon-margin-right":N}}),s=b?M("alert",R(()=>r.type[0]),g,r):void 0,h=Z(!0),C=()=>{const{onAfterLeave:n,onAfterHide:e}=r;n&&n(),e&&e()};return{rtlEnabled:y,mergedClsPrefix:o,mergedBordered:d,visible:h,handleCloseClick:()=>{var n;Promise.resolve((n=r.onClose)===null||n===void 0?void 0:n.call(r)).then(e=>{e!==!1&&(h.value=!1)})},handleAfterLeave:()=>{C()},mergedTheme:u,cssVars:b?void 0:g,themeClass:s==null?void 0:s.themeClass,onRender:s==null?void 0:s.onRender}},render(){var r;return(r=this.onRender)===null||r===void 0||r.call(this),t(eo,{onAfterLeave:this.handleAfterLeave},{default:()=>{const{mergedClsPrefix:o,$slots:d}=this,b={class:[`${o}-alert`,this.themeClass,this.closable&&`${o}-alert--closable`,this.showIcon&&`${o}-alert--show-icon`,!this.title&&this.closable&&`${o}-alert--right-adjust`,this.rtlEnabled&&`${o}-alert--rtl`],style:this.cssVars,role:"alert"};return this.visible?t("div",Object.assign({},oo(this.$attrs,b)),this.closable&&t(V,{clsPrefix:o,class:`${o}-alert__close`,onClick:this.handleCloseClick}),this.bordered&&t("div",{class:`${o}-alert__border`}),this.showIcon&&t("div",{class:`${o}-alert__icon`,"aria-hidden":"true"},J(d.icon,()=>[t(D,{clsPrefix:o},{default:()=>{switch(this.type){case"success":return t(to,null);case"info":return t(so,null);case"warning":return t(no,null);case"error":return t(ro,null);default:return null}}})])),t("div",{class:[`${o}-alert-body`,this.mergedBordered&&`${o}-alert-body--bordered`]},K(d.header,l=>{const u=l||this.title;return u?t("div",{class:`${o}-alert-body__title`},u):null}),d.default&&t("div",{class:`${o}-alert-body__content`},d))):null}})}});export{Co as N};
