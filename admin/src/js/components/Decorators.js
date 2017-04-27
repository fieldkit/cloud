/* @flow */
 import React, { Component } from 'react'
 import { SketchPicker } from 'react-color';
 import ColorBrewer from 'colorbrewer';
  import type { Lens, Lens_ } from 'safety-lens'
  import { get, set, compose } from 'safety-lens'
  import { prop, _1, _2 } from 'safety-lens/es2015'
 import type {Attr} from './Collection';
  
  export type Stop = {
      location: number;
     color: string;
 }
 
 export type InterpolationType = "constant" | "linear"
 
 export type Color = {
     type: InterpolationType;
     colors: Stop[];
     data_key: ?string;
     bounds: ?[number,number];
 }
 
 const _colorType: Lens_<Color,InterpolationType> = prop("type")
 const _colorColors: Lens_<Color,Stop[]> = prop("colors")
 const _colorDataKey: Lens_<Color,?string> = prop("data_key")
 const _colorBounds: Lens_<Color,?[number,number]> = prop("bounds")
 
 export type Size = {
     type: InterpolationType;
     data_key: ?string;
     bounds: [number,number];
 }
 
 const _sizeType: Lens_<Size,InterpolationType> = prop("type")
 const _sizeDataKey: Lens_<Size,?string> = prop("data_key")
 const _sizeBounds: Lens_<Size,[number,number]> = prop("bounds")
 
 export type PointDecorator = {
     collection_id: string;
     points: {
         color: Color,
         size: Size,
         sprite: string
     };
     title: string;
     type: "point";
 }
 
 const _pointDecoratorPointsColor: Lens_<PointDecorator,Color> = compose(prop("points"),prop("color"))
 const _pointDecoratorPointsSize: Lens_<PointDecorator,Size> = compose(prop("points"),prop("size"))
 const _pointDecoratorPointsSprite: Lens_<PointDecorator,string> = compose(prop("points"),prop("sprite"))
 const _pointDecoratorTitle: Lens_<PointDecorator,string> = prop("title")
 
 export function updatePointDecorator<A>(l: Lens_<PointDecorator,A>,value:A,p:PointDecorator):PointDecorator{
     return set(l,value,p)
 }
 
 export function emptyPointDecorator(): PointDecorator{
     return {
         collection_id: "",
         points: {
             color: {
                 type: "constant",
                 colors: [{location: 0, color: "#ff0000"}],
                 data_key: null,
                 bounds: null
             },
             size: {
                 type: "constant",
                 data_key: null,
                 bounds: [15,15]
             },
             sprite: "circle.png"
         },
         title: "",
         type: "point"
     }
 }
  
  export type Decorator = PointDecorator
  type PointDecoratorProps = {
     initial_state: PointDecorator;
     attributes: {[string]: Attr};
  }
  
  
  export class PointDecoratorComponent extends Component {
  props: {initial_state: PointDecorator}
     props: PointDecoratorProps
      state: {data: PointDecorator, schemer_open: boolean}
      toggleColorType: () => void
      toggleSizeType: () => void
      setLowerSize: (Object) => void
      setUpperSize: (Object) => void
     setSize: (Object) => void
      setSprite: (Object) => void
     updateSizeDataKey: (Object) => void
     updateColorDataKey: (Object) => void
  
      constructor(props: PointDecoratorProps){
          super(props)
          this.state = {data: this.props.initial_state, schemer_open: false}
          this.toggleColorType = this.toggleColorType.bind(this)
          this.toggleSizeType = this.toggleSizeType.bind(this)
          this.setSprite = this.setSprite.bind(this)
         this.updateSizeDataKey = this.updateSizeDataKey.bind(this)
         this.updateColorDataKey = this.updateColorDataKey.bind(this)
      }
   
   setLowerSize(e:Object){
 
     update<A>(lens: Lens_<PointDecorator,A>,value:A): void{
          let {data} = this.state;
         data = updatePointDecorator(lens,value,data)
         this.setState({data})
     }
   
     setSize(e:Object){
         const value = Number(e.target.value);
         const size_lens = compose(_pointDecoratorPointsSize,_sizeBounds)
         this.update(size_lens,[value,value])
     }
 
     setLowerSize(e:Object){
          const value = Number(e.target.value);
          const size_lens = compose(_pointDecoratorPointsSize,_sizeBounds,_1)
       data = updatePointDecorator(size_lens,value,data)
       this.setState({data})
         this.update(size_lens,value)
      }
      
      setUpperSize(e:Object){
      let {data} = this.state;
          const value = Number(e.target.value);
          const size_lens = compose(_pointDecoratorPointsSize,_sizeBounds,_2)
      data = updatePointDecorator(size_lens,value,data)
      this.setState({data})
         this.update(size_lens,value)
      }
  
      setSprite(e: Object){
      let {data} = this.state;
          const value = e.target.value;
      data = updatePointDecorator(_pointDecoratorPointsSprite,value,data)
      this.setState({data})
         this.update(_pointDecoratorPointsSprite,value)
      }
  
      toggleColorType(){
         const {data} = this.state;
         if(data.points.color.type === "constant"){
             this.setBrewerColors(ColorBrewer.Reds[5]) 
         } else {
             this.setConstantColor(data.points.color.colors[0].color) 
         }
     }
     
     toggleSizeType(){
          let {data} = this.state;
          let size_lens = compose(_pointDecoratorPointsSize,_sizeType)
          if(data.points.size.type === "constant"){
          data = updatePointDecorator(size_lens,"linear",data)
             this.update(size_lens,"linear")
          } else {
          data = updatePointDecorator(size_lens,"constant",data)
             this.update(size_lens,"constant")
          }
      this.setState({data})
      }
  
      setBrewerColors(brewer_colors: string[]){
         let color_length = brewer_colors.length || 0,
             step: number, 
             colors: Stop[];
             
         if(color_length > 1){
             step = 1/(color_length - 1)
             colors = brewer_colors.map((c,i) => { 
                 let location
                 if(i < color_length - 1){
                     location = i*step;
                 } else {
                     location = 1.0
                 }
                 return {location, color: c}
             })  
             this.setLinearColor(colors)
         } else if (color_length === 1){
             this.setConstantColor(brewer_colors[0])
         }
     }
 
     setLinearColor(colors: Stop[]){
             let color_lens = compose(_pointDecoratorPointsColor,_colorColors)
             let color_type_lens = compose(_pointDecoratorPointsColor,_colorType)
             let data = updatePointDecorator(color_lens,colors,this.state.data)
             data = updatePointDecorator(color_type_lens,"linear",data)
             this.setState({data, schemer_open: false})
     }
 
     setConstantColor(color: string){
         let new_color = [{location: 0, color: color}]
         let color_lens = compose(_pointDecoratorPointsColor,_colorColors)
         let color_type_lens = compose(_pointDecoratorPointsColor,_colorType)
         let data = updatePointDecorator(color_lens,new_color,this.state.data)
         data = updatePointDecorator(color_type_lens,"constant",data)
          this.setState({data, schemer_open:false})
      }
  
     updateSizeDataKey(e: Object){
         const value = e.target.value;
         const size_lens = compose(_pointDecoratorPointsSize,_sizeDataKey)
         this.update(size_lens,value)
     }
 
     updateColorDataKey(e: Object){
         const value = e.target.value;
         const color_lens = compose(_pointDecoratorPointsColor,_colorDataKey)
         this.update(color_lens,value)
     }
 
      render(){
         const {data} = this.state;
        const {attributes} = this.props;
         const collections = null
        const target_attrs = Object.keys(attributes).filter(a => attributes[a].type === "num")
                                   .map((a) => {
                                     const attr = attributes[a];                 
                                     return (
                                         <option value="{attr.name}">{attr.name}</option>
                                     )
                                 })
         let color,size; 
  
         if( data.points.color.type === "constant"){
           color = <SketchPicker onChangeComplete={c => this.setConstantColor(c.hex)} color={data.points.color.colors[0].color} disableAlpha={true}/>
         } else {
       let brewer_selections = Object.keys(ColorBrewer).map((k) => {
          let scheme = ColorBrewer[k][5]
          const brewer_selections = Object.keys(ColorBrewer).map((k) => {
             const scheme = ColorBrewer[k][5]
              return (
                  <div className="brewer-selection" key={k} onClick={() => this.setBrewerColors(scheme)}>
                      {scheme.map((hex,i) => {
                         return (
                             <div className='color-thumb' style={{backgroundColor:hex}} key={i}></div>
                         )
                     })}
                  </div>
              )
           })
          
 
           color = (
              <div>
                 <h4>Based on:</h4>
                 <select value={data.points.color.data_key} onChange={this.updateColorDataKey}>
                     {target_attrs}
                 </select>
                 <h4>Color: </h4>
                  {data.points.color.colors.map((c,i) => {
                      return (
                          <div className='color-thumb' style={{backgroundColor:c.color}} key={i}></div>
                     )
                 })}
                 <div>
                     <button className="brewer-title" onClick={() => this.setState({schemer_open: true})}>Select Color Scheme</button>
                     <div className="brewer-schemes" style={{display: this.state.schemer_open ? "block" : "none"}}>
                         {brewer_selections}
                     </div>
                 </div>
             </div>
          )
         }
  
         if(data.points.size.type === "constant"){
         size = <input value={data.points.size.bounds[0]} onChange={e => this.setLowerSize(e)}/>
          size = <input value={data.points.size.bounds[0]} onChange={e => this.setSize(e)}/>
         } else {
           size = (
               <div>
                 <h4>Based on</h4>
                 <select value={data.points.size.data_key} onChange={this.updateSizeDataKey}>
                     {target_attrs}
                 </select>
                 <h4>Max/min</h4>
                   <input value={data.points.size.bounds[0]} onChange={e => this.setLowerSize(e)}/>
                   <input value={data.points.size.bounds[1]} onChange={e => this.setUpperSize(e)}/>
               </div>
          )
        }
         
        return (
            <div className="point-decorator">
                 <div className="decorator-row">
                     <div className="decorator-row-label">Source: </div>
                     <select>
                         {collections}
                     </select>
                 </div>
                 <div className="decorator-row">
                     <span className="decorator-row-label">Color: </span>
                     <select value={data.points.color.type} onChange={this.toggleColorType}>
                         <option value="constant">constant</option>
                         <option value="linear">linear</option>
                     </select>
                     {color}
                 </div>
                 <div className="decorator-row">
                     <span className="decorator-row-label">Size: </span>
                     <select value={data.points.size.type} onChange={this.toggleSizeType}>
                         <option value="constant">constant</option>
                         <option value="linear">linear</option>
                     </select>
                     {size}
                 </div>
                 <div className="decorator-row">
                     <span className="decorator-row-label">Sprite: </span>
                     <input value={data.points.sprite} onChange={this.setSprite}/>
                 </div>
            </div>
        )
     }
 }