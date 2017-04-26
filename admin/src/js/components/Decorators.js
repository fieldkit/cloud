/* @flow */
import React, { Component } from 'react'
import { SketchPicker } from 'react-color';
import ColorBrewer from 'colorbrewer';
import type { Lens, Lens_ } from 'safety-lens'
import { get, set, compose } from 'safety-lens'
import { prop, _1, _2 } from 'safety-lens/es2015'

import Dropdown, { DropdownTrigger, DropdownContent } from 'react-simple-dropdown';
import { FormItem } from './forms/FormItem'
import { FormSelectItem } from './forms/FormSelectItem'
import type { APIErrors } from '../api/types';
import '../../css/decorators.css'

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
  intial_state: PointDecorator
}


export class PointDecoratorComponent extends Component {
  props: {
    initial_state: PointDecorator
  }
  state: {
    data: PointDecorator,
    errors: ?APIErrors
  }
  toggleColorType: () => void
  toggleSizeType: () => void
  setLowerSize: (Object) => void
  setUpperSize: (Object) => void
  setSprite: (Object) => void

  constructor(props: PointDecoratorProps){
    super(props)
    this.state = {
      data: this.props.initial_state,
      errors: null
    }
    this.toggleColorType = this.toggleColorType.bind(this)
    this.toggleSizeType = this.toggleSizeType.bind(this)
    this.setSprite = this.setSprite.bind(this)
  }
   
  setLowerSize(e:Object){
    let {data} = this.state;
    const value = Number(e.target.value);
    const size_lens = compose(_pointDecoratorPointsSize,_sizeBounds,_1)
    data = updatePointDecorator(size_lens,value,data)
    this.setState({data})
  }
  
  setUpperSize(e:Object){
    let {data} = this.state;
    const value = Number(e.target.value);
    const size_lens = compose(_pointDecoratorPointsSize,_sizeBounds,_2)
    data = updatePointDecorator(size_lens,value,data)
    this.setState({data})
  }

  setSprite(e: Object){
    let {data} = this.state;
    const value = e.target.value;
    data = updatePointDecorator(_pointDecoratorPointsSprite,value,data)
    this.setState({data})
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
    } else {
      data = updatePointDecorator(size_lens,"constant",data)
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
    this.setState({data})
  }

  setConstantColor(color: string){
    let new_color = [{location: 0, color: color}]
    let color_lens = compose(_pointDecoratorPointsColor,_colorColors)
    let color_type_lens = compose(_pointDecoratorPointsColor,_colorType)
    let data = updatePointDecorator(color_lens,new_color,this.state.data)
    data = updatePointDecorator(color_type_lens,"constant",data)
    this.setState({data})
  }

  render(){
    const { data, errors } = this.state;
    const options = [{value: 'constant', text: 'constant'}, {value: 'linear', text: 'linear'}];
    const collections = null
    let colorDropdownTrigger, colorDropdownContent, size; 
    const divStyle = {
      backgroundColor: data.points.color.colors[0].color
    }

    if( data.points.color.type === "constant"){
    
      colorDropdownTrigger = <div className="color-form" style={divStyle}></div>;
      colorDropdownContent = <SketchPicker onChangeComplete={c => this.setConstantColor(c.hex)} color={data.points.color.colors[0].color} disableAlpha={true}/>;
    
    } else {

      let brewer_selections = Object.keys(ColorBrewer).map((k) => {
        let scheme = ColorBrewer[k][5]
        return (
          <div className="brewer-selection" key={k} onClick={() => this.setBrewerColors(scheme)}>
            {scheme.map((hex,i) => {
              return (
                <div className='color-thumb' style={{backgroundColor:hex}} key={i}></div>
              )
            })}
          </div>
        )
      });

      colorDropdownTrigger = (
        <div>
          {data.points.color.colors.map((c,i) => {
              return (
                <div className='color-thumb' style={{backgroundColor:c.color}} key={i}></div>
              )
            })
          }
        </div>
      );

      colorDropdownContent = (
        <div>
          <div className="brewer-schemes">
            {brewer_selections}
          </div>
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
          <FormSelectItem
            labelText={'Color'}
            name={'color'}
            value={data.points.color.type}
            inline={true}
            firstOptionText={'Select'}
            options={options}
            errors={errors}
            onChange={this.toggleColorType}
          />
          <Dropdown className="color-dropdown" ref="color-dropdown">
            <DropdownTrigger className="trigger">
              { colorDropdownTrigger }
            </DropdownTrigger>
            <DropdownContent className="dropdown-contents">
              { colorDropdownContent }
            </DropdownContent>
          </Dropdown>          
        </div>
        <div className="decorator-row">
          <FormSelectItem
            labelText={'Size'}
            name={'size'}
            value={data.points.size.type}
            inline={true}
            firstOptionText={'Select'}
            options={options}
            errors={errors}
            onChange={this.toggleSizeType}
          />
          <FormItem
            labelText={'Value'}
            name={'value'}
            value={data.points.size.bounds[0]}
            inline={true}
            errors={errors}
            onChange={e => this.setLowerSize(e)}
          />
          { data.points.size.type !== "constant" &&
            <FormItem
              labelText={'Max'}
              name={'upper-size'}
              value={data.points.size.bounds[1]}
              inline={true}
              errors={errors}
              onChange={e => this.setUpperSize(e)}
            />          
          }
        </div>
        <div className="decorator-row">
          <FormItem
            labelText={'Sprite'}
            name={'sprite'}
            value={data.points.sprite}
            errors={errors}
            onChange={this.setSprite}
          />
        </div>
      </div>
    )
  }
}
