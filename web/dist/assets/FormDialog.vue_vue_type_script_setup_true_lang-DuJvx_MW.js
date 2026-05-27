import{N as ge}from"./Modal-rM18g_Zj.js";import{u as ye,B as G}from"./Button-Bkv4AwgB.js";import{N as we}from"./InputNumber-cWCfWutJ.js";import{b as ke}from"./Select-BSIolqJY.js";import{p as J,s as H,C as F,c as xe,b as Ce,F as Se,L as Re,e as X,x as q,g as y}from"./Scrollbar-DUkKw_1i.js";import{a5 as $e,I as Be,z as ze,t as Z,v as l,s as Y,w,x as ee,a3 as ae,al as c,by as oe,b4 as j,N as P,X as N,bk as Ve,bA as Fe,Q as k,bC as _,aO as Ne,bo as h,aV as m,a0 as W,V as _e,b5 as Ue,F as De,_ as te,R as je}from"./index-qenDAa9H.js";import{u as Oe}from"./get-K_rslPNJ.js";import{u as Te}from"./use-message-BcbmdW1g.js";import{N as Pe,a as We}from"./FormItem-BUTwc8nL.js";import{N as Me}from"./Space-D6xIht3K.js";import{a as M}from"./Input-8oJ3ABXT.js";function Ie(e){const{primaryColor:b,opacityDisabled:f,borderRadius:s,textColor3:p}=e;return Object.assign(Object.assign({},Be),{iconColor:p,textColor:"white",loadingColor:b,opacityDisabled:f,railColor:"rgba(0, 0, 0, .14)",railColorActive:b,buttonBoxShadow:"0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)",buttonColor:"#FFF",railBorderRadiusSmall:s,railBorderRadiusMedium:s,railBorderRadiusLarge:s,buttonBorderRadiusSmall:s,buttonBorderRadiusMedium:s,buttonBorderRadiusLarge:s,boxShadowFocus:`0 0 0 2px ${ze(b,{alpha:.2})}`})}const Ke={common:$e,self:Ie},Ae=Z("switch",`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[l("children-placeholder",`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),l("rail-placeholder",`
 display: flex;
 flex-wrap: none;
 `),l("button-placeholder",`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),Z("base-loading",`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[J({left:"50%",top:"50%",originalTransform:"translateX(-50%) translateY(-50%)"})]),l("checked, unchecked",`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 box-sizing: border-box;
 position: absolute;
 white-space: nowrap;
 top: 0;
 bottom: 0;
 display: flex;
 align-items: center;
 line-height: 1;
 `),l("checked",`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),l("unchecked",`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),Y("&:focus",[l("rail",`
 box-shadow: var(--n-box-shadow-focus);
 `)]),w("round",[l("rail","border-radius: calc(var(--n-rail-height) / 2);",[l("button","border-radius: calc(var(--n-button-height) / 2);")])]),ee("disabled",[ee("icon",[w("rubber-band",[w("pressed",[l("rail",[l("button","max-width: var(--n-button-width-pressed);")])]),l("rail",[Y("&:active",[l("button","max-width: var(--n-button-width-pressed);")])]),w("active",[w("pressed",[l("rail",[l("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])]),l("rail",[Y("&:active",[l("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])])])])])]),w("active",[l("rail",[l("button","left: calc(100% - var(--n-button-width) - var(--n-offset))")])]),l("rail",`
 overflow: hidden;
 height: var(--n-rail-height);
 min-width: var(--n-rail-width);
 border-radius: var(--n-rail-border-radius);
 cursor: pointer;
 position: relative;
 transition:
 opacity .3s var(--n-bezier),
 background .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-rail-color);
 `,[l("button-icon",`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 font-size: calc(var(--n-button-height) - 4px);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 display: flex;
 justify-content: center;
 align-items: center;
 line-height: 1;
 `,[J()]),l("button",`
 align-items: center; 
 top: var(--n-offset);
 left: var(--n-offset);
 height: var(--n-button-height);
 width: var(--n-button-width-pressed);
 max-width: var(--n-button-width);
 border-radius: var(--n-button-border-radius);
 background-color: var(--n-button-color);
 box-shadow: var(--n-button-box-shadow);
 box-sizing: border-box;
 cursor: inherit;
 content: "";
 position: absolute;
 transition:
 background-color .3s var(--n-bezier),
 left .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 `)]),w("active",[l("rail","background-color: var(--n-rail-color-active);")]),w("loading",[l("rail",`
 cursor: wait;
 `)]),w("disabled",[l("rail",`
 cursor: not-allowed;
 opacity: .5;
 `)])]),Ee=Object.assign(Object.assign({},oe.props),{size:String,value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},spinProps:Object,onChange:[Function,Array]});let T;const Le=ae({name:"Switch",props:Ee,slots:Object,setup(e){T===void 0&&(typeof CSS<"u"?typeof CSS.supports<"u"?T=CSS.supports("width","max(1px)"):T=!1:T=!0);const{mergedClsPrefixRef:b,inlineThemeDisabled:f,mergedComponentPropsRef:s}=Se(e),p=oe("Switch","-switch",Ae,Ke,e,b),d=ye(e,{mergedSize(a){var S,R;if(e.size!==void 0)return e.size;if(a)return a.mergedSize.value;const V=(R=(S=s==null?void 0:s.value)===null||S===void 0?void 0:S.Switch)===null||R===void 0?void 0:R.size;return V||"medium"}}),{mergedSizeRef:x,mergedDisabledRef:o}=d,z=j(e.defaultValue),U=Ve(e,"value"),g=Oe(U,z),C=P(()=>g.value===e.checkedValue),r=j(!1),v=j(!1),i=P(()=>{const{railStyle:a}=e;if(a)return a({focused:v.value,checked:C.value})});function n(a){const{"onUpdate:value":S,onChange:R,onUpdateValue:V}=e,{nTriggerFormInput:I,nTriggerFormChange:K}=d;S&&X(S,a),V&&X(V,a),R&&X(R,a),z.value=a,I(),K()}function t(){const{nTriggerFormFocus:a}=d;a()}function u(){const{nTriggerFormBlur:a}=d;a()}function le(){e.loading||o.value||(g.value!==e.checkedValue?n(e.checkedValue):n(e.uncheckedValue))}function ne(){v.value=!0,t()}function re(){v.value=!1,u(),r.value=!1}function ie(a){e.loading||o.value||a.key===" "&&(g.value!==e.checkedValue?n(e.checkedValue):n(e.uncheckedValue),r.value=!1)}function se(a){e.loading||o.value||a.key===" "&&(a.preventDefault(),r.value=!0)}const Q=P(()=>{const{value:a}=x,{self:{opacityDisabled:S,railColor:R,railColorActive:V,buttonBoxShadow:I,buttonColor:K,boxShadowFocus:ue,loadingColor:ce,textColor:de,iconColor:he,[N("buttonHeight",a)]:$,[N("buttonWidth",a)]:ve,[N("buttonWidthPressed",a)]:be,[N("railHeight",a)]:B,[N("railWidth",a)]:O,[N("railBorderRadius",a)]:pe,[N("buttonBorderRadius",a)]:me},common:{cubicBezierEaseInOut:fe}}=p.value;let A,E,L;return T?(A=`calc((${B} - ${$}) / 2)`,E=`max(${B}, ${$})`,L=`max(${O}, calc(${O} + ${$} - ${B}))`):(A=q((y(B)-y($))/2),E=q(Math.max(y(B),y($))),L=y(B)>y($)?O:q(y(O)+y($)-y(B))),{"--n-bezier":fe,"--n-button-border-radius":me,"--n-button-box-shadow":I,"--n-button-color":K,"--n-button-width":ve,"--n-button-width-pressed":be,"--n-button-height":$,"--n-height":E,"--n-offset":A,"--n-opacity-disabled":S,"--n-rail-border-radius":pe,"--n-rail-color":R,"--n-rail-color-active":V,"--n-rail-height":B,"--n-rail-width":O,"--n-width":L,"--n-box-shadow-focus":ue,"--n-loading-color":ce,"--n-text-color":de,"--n-icon-color":he}}),D=f?Re("switch",P(()=>x.value[0]),Q,e):void 0;return{handleClick:le,handleBlur:re,handleFocus:ne,handleKeyup:ie,handleKeydown:se,mergedRailStyle:i,pressed:r,mergedClsPrefix:b,mergedValue:g,checked:C,mergedDisabled:o,cssVars:f?void 0:Q,themeClass:D==null?void 0:D.themeClass,onRender:D==null?void 0:D.onRender}},render(){const{mergedClsPrefix:e,mergedDisabled:b,checked:f,mergedRailStyle:s,onRender:p,$slots:d}=this;p==null||p();const{checked:x,unchecked:o,icon:z,"checked-icon":U,"unchecked-icon":g}=d,C=!(H(z)&&H(U)&&H(g));return c("div",{role:"switch","aria-checked":f,class:[`${e}-switch`,this.themeClass,C&&`${e}-switch--icon`,f&&`${e}-switch--active`,b&&`${e}-switch--disabled`,this.round&&`${e}-switch--round`,this.loading&&`${e}-switch--loading`,this.pressed&&`${e}-switch--pressed`,this.rubberBand&&`${e}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},c("div",{class:`${e}-switch__rail`,"aria-hidden":"true",style:s},F(x,r=>F(o,v=>r||v?c("div",{"aria-hidden":!0,class:`${e}-switch__children-placeholder`},c("div",{class:`${e}-switch__rail-placeholder`},c("div",{class:`${e}-switch__button-placeholder`}),r),c("div",{class:`${e}-switch__rail-placeholder`},c("div",{class:`${e}-switch__button-placeholder`}),v)):null)),c("div",{class:`${e}-switch__button`},F(z,r=>F(U,v=>F(g,i=>c(xe,null,{default:()=>this.loading?c(Ce,Object.assign({key:"loading",clsPrefix:e,strokeWidth:20},this.spinProps)):this.checked&&(v||r)?c("div",{class:`${e}-switch__button-icon`,key:v?"checked-icon":"icon"},v||r):!this.checked&&(i||r)?c("div",{class:`${e}-switch__button-icon`,key:i?"unchecked-icon":"icon"},i||r):null})))),F(x,r=>r&&c("div",{key:"checked",class:`${e}-switch__checked`},r)),F(o,r=>r&&c("div",{key:"unchecked",class:`${e}-switch__unchecked`},r)))))}}),ot=ae({__name:"FormDialog",props:{title:{},fields:{},initialData:{},loading:{type:Boolean},width:{},labelWidth:{}},emits:["submit","cancel","update:show"],setup(e,{expose:b,emit:f}){const s=e,p=f,d=j(!1),x=j(null),o=j({}),z=Te(),U=P(()=>{const i={};return s.fields.forEach(n=>{n.required?i[n.key]=[...n.rules||[],{required:!0,message:`请输入${n.label}`,trigger:["blur","change"]}]:n.rules&&(i[n.key]=n.rules)}),i});function g(i){d.value=!0,r(i||s.initialData||{})}function C(){d.value=!1,p("cancel")}function r(i){const n={};s.fields.forEach(t=>{n[t.key]=i[t.key]!==void 0?i[t.key]:t.defaultValue??""}),o.value=n}async function v(){var i;try{await((i=x.value)==null?void 0:i.validate()),p("submit",{...o.value})}catch{z.warning("请检查表单填写")}}return Fe(()=>s.initialData,i=>{i&&d.value&&r(i)},{deep:!0}),b({open:g,close:C}),(i,n)=>(m(),k(h(ge),{show:d.value,"onUpdate:show":n[0]||(n[0]=t=>d.value=t),title:e.title,preset:"card",style:Ne({width:`${e.width||600}px`}),"mask-closable":!1,onClose:C},{footer:_(()=>[W(h(Me),{justify:"end"},{default:_(()=>[W(h(G),{onClick:C},{default:_(()=>[...n[1]||(n[1]=[te("取消",-1)])]),_:1}),W(h(G),{type:"primary",loading:e.loading,onClick:v},{default:_(()=>[...n[2]||(n[2]=[te(" 确定 ",-1)])]),_:1},8,["loading"])]),_:1})]),default:_(()=>[W(h(Pe),{ref_key:"formRef",ref:x,model:o.value,rules:U.value,"label-width":e.labelWidth||100,"label-placement":"left"},{default:_(()=>[(m(!0),_e(De,null,Ue(e.fields,t=>(m(),k(h(We),{key:t.key,label:t.label,path:t.key},{default:_(()=>[t.type==="text"?(m(),k(h(M),{key:0,value:o.value[t.key],"onUpdate:value":u=>o.value[t.key]=u,placeholder:t.placeholder||`请输入${t.label}`,disabled:t.disabled},null,8,["value","onUpdate:value","placeholder","disabled"])):t.type==="password"?(m(),k(h(M),{key:1,value:o.value[t.key],"onUpdate:value":u=>o.value[t.key]=u,type:"password","show-password-on":"click",placeholder:t.placeholder||`请输入${t.label}`},null,8,["value","onUpdate:value","placeholder"])):t.type==="textarea"?(m(),k(h(M),{key:2,value:o.value[t.key],"onUpdate:value":u=>o.value[t.key]=u,type:"textarea",placeholder:t.placeholder||`请输入${t.label}`,autosize:{minRows:3,maxRows:8}},null,8,["value","onUpdate:value","placeholder"])):t.type==="code"?(m(),k(h(M),{key:3,value:o.value[t.key],"onUpdate:value":u=>o.value[t.key]=u,type:"textarea",placeholder:t.placeholder||`请输入${t.label}`,autosize:{minRows:5,maxRows:15},style:{"font-family":"monospace"}},null,8,["value","onUpdate:value","placeholder"])):t.type==="number"?(m(),k(h(we),{key:4,value:o.value[t.key],"onUpdate:value":u=>o.value[t.key]=u,placeholder:t.placeholder||`请输入${t.label}`,min:t.min,max:t.max},null,8,["value","onUpdate:value","placeholder","min","max"])):t.type==="select"?(m(),k(h(ke),{key:5,value:o.value[t.key],"onUpdate:value":u=>o.value[t.key]=u,options:t.options||[],placeholder:t.placeholder||`请选择${t.label}`},null,8,["value","onUpdate:value","options","placeholder"])):t.type==="switch"?(m(),k(h(Le),{key:6,value:o.value[t.key],"onUpdate:value":u=>o.value[t.key]=u},null,8,["value","onUpdate:value"])):je("",!0)]),_:2},1032,["label","path"]))),128))]),_:1},8,["model","rules","label-width"])]),_:1},8,["show","title","style"]))}});export{Le as N,ot as _};
