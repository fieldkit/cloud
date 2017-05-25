/* @flow */
import React, { Component } from 'react'
import { SketchPicker } from 'react-color';
import ReactModal from 'react-modal';
import ColorBrewer from 'colorbrewer';
import Dropdown, { DropdownTrigger, DropdownContent } from 'react-simple-dropdown';
import type { Lens_ } from 'safety-lens';
import { set, compose } from 'safety-lens';
import { prop, _1, _2 } from 'safety-lens/es2015';
import { FormItem } from '../forms/FormItem';
import { FormSelectItem } from '../forms/FormSelectItem';
import type { APIErrors } from '../../api/types';
import type { ProjectData } from '../../types/CollectionTypes';
import type { Stop, InterpolationType, Color, Size, PointDecorator, Decorator, Viz } from '../../types/VizTypes';

import '../../../css/decorators.css';

const _colorType: Lens_<Color,InterpolationType> = prop('type');
const _colorColors: Lens_<Color,Stop[]> = prop('colors');
const _colorDataKey: Lens_<Color,?string> = prop('data_key');
const _colorBounds: Lens_<Color,?[number,number]> = prop('bounds');

const _sizeType: Lens_<Size,InterpolationType> = prop('type');
const _sizeDataKey: Lens_<Size,?string> = prop('data_key');
const _sizeBounds: Lens_<Size,[number,number]> = prop('bounds');

const _pointDecoratorPointsLocation: Lens_<PointDecorator,string> = compose(prop('points'),prop('location'));
const _pointDecoratorPointsColor: Lens_<PointDecorator,Color> = compose(prop('points'),prop('color'));
const _pointDecoratorPointsSize: Lens_<PointDecorator,Size> = compose(prop('points'),prop('size'));
const _pointDecoratorPointsSprite: Lens_<PointDecorator,string> = compose(prop('points'),prop('sprite'));
const _pointDecoratorTitle: Lens_<PointDecorator,string> = prop('title');

export function updatePointDecorator<A>(l: Lens_<PointDecorator,A>,value:A,p:PointDecorator):PointDecorator{
  return set(l,value,p);
}

export const _decorator: Lens_<Viz,Decorator> = prop('decorator');

type PointDecoratorProps = {
   initial_state: PointDecorator;
   project_data: ProjectData;
}

export default class PointDecoratorComponent extends Component {
  state: {
    errors: ?APIErrors
  }
  props: {
    viz: Viz,
    creator: Object
  }
  toggleColorType: () => void
  toggleSizeType: () => void
  setLowerSize: (Object) => void
  setUpperSize: (Object) => void
  setSize: (Object) => void
  setLocationKey: (Object) => void
  setSprite: (Object) => void
  updateSizeDataKey: (Object) => void
  updateColorDataKey: (Object) => void
  save: () => void

  constructor(props: PointDecoratorProps){
    super(props)

    this.toggleColorType = this.toggleColorType.bind(this)
    this.toggleSizeType = this.toggleSizeType.bind(this)
    this.setSprite = this.setSprite.bind(this)
    this.setSize = this.setSize.bind(this)
    this.setLocationKey = this.setLocationKey.bind(this)
    this.updateSizeDataKey = this.updateSizeDataKey.bind(this)
    this.updateColorDataKey = this.updateColorDataKey.bind(this)
    this.save = this.save.bind(this)
    this.state = {
      errors: null
    }
  }

  update<A>(lens: Lens_<PointDecorator,A>,value:A): void{
    let data = this.props.viz.decorator;
    data = updatePointDecorator(lens,value,data)
    this.props.creator.update(_decorator,data)
  }

  setLocationKey(e:Object){
    const value = String(e.target.value);
    this.update(_pointDecoratorPointsLocation,value)
  }

  setSize(e:Object){
    const value = Number(e.target.value);
    const size_lens = compose(_pointDecoratorPointsSize,_sizeBounds)
    this.update(size_lens,[value,value])
  }

  setLowerSize(e:Object){
    const value = Number(e.target.value);
    const size_lens = compose(_pointDecoratorPointsSize,_sizeBounds,_1)
    this.update(size_lens,value)
  }

  setUpperSize(e:Object){
    const value = Number(e.target.value);
    const size_lens = compose(_pointDecoratorPointsSize,_sizeBounds,_2)
    this.update(size_lens,value)
  }

  setSprite(e: Object){
    const value = e.target.value;
    this.update(_pointDecoratorPointsSprite,value)
  }

  save(){
    let serialized_viz = JSON.stringify(this.props.viz)
    console.log(serialized_viz)
  }

  toggleColorType(){
    const data = this.props.viz.decorator;
    if(data.points.color.type === 'constant'){
      this.setBrewerColors(ColorBrewer.Reds[5])
    } else {
      this.setConstantColor(data.points.color.colors[0].color)
    }
  }

  toggleSizeType(){
    const data = this.props.viz.decorator;
    let size_lens = compose(_pointDecoratorPointsSize,_sizeType)
    if(data.points.size.type === 'constant'){
      this.update(size_lens,'linear')
    } else {
      this.update(size_lens,'constant')
    }
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
    let data = this.props.viz.decorator;
    data = updatePointDecorator(color_lens,colors,data)
    data = updatePointDecorator(color_type_lens,'linear',data)
    this.props.creator.update(_decorator,data)
  }

  setConstantColor(color: string){
    let new_color = [{location: 0, color: color}]
    let color_lens = compose(_pointDecoratorPointsColor,_colorColors)
    let color_type_lens = compose(_pointDecoratorPointsColor,_colorType)
    let data = this.props.viz.decorator;
    data = updatePointDecorator(color_lens,new_color,data)
    data = updatePointDecorator(color_type_lens,'constant',data)
    this.props.creator.update(_decorator,data)
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
    const { errors } = this.state;
    const {viz} = this.props
    const data = viz.decorator;
    const target_attrs = viz.selection_operations.map((s) => {
      return {value: s.id, text: s.value_name}
    })
    let options;
    if(viz.selection_operations.length > 0){
      options = [{value: 'constant', text: 'constant'}, {value: 'linear', text: 'linear'}];
    } else {
      options = [{value: 'constant', text: 'constant'}];
    }
    const collections = null
    let colorDropdownTrigger, colorDropdownContent, size;

    if( data.points.color.type === 'constant'){

      colorDropdownTrigger = (
        <div className='selected-color'>
          <div className='color-thumb' style={{ backgroundColor: data.points.color.colors[0].color }}></div>
        </div>
      )
      colorDropdownContent = <SketchPicker onChangeComplete={c => this.setConstantColor(c.hex)} color={data.points.color.colors[0].color} disableAlpha={true}/>;

    } else {


      const brewer_selections = Object.keys(ColorBrewer).map((k) => {
        const scheme = ColorBrewer[k][5]
        return (
          <div className='brewer-selection' key={k} onClick={() => this.setBrewerColors(scheme)}>
            {scheme.map((hex,i) => {
              return (
                <div className='color-thumb' style={{backgroundColor:hex}} key={i}></div>
              )
            })}
          </div>
        )
      });

      colorDropdownTrigger = (
        <div className='selected-color'>
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
          <div className='brewer-schemes'>
            {brewer_selections}
          </div>
        </div>
      )
    }

    return (
      <div className='point-decorator'>

        <div className='decorator-row'>
          <FormSelectItem
            labelText={'Location'}
            name={'location'}
            value={data.points.location}
            inline={false}
            firstOptionText={'Select'}
            options={target_attrs}
            errors={errors}
            onChange={this.setLocationKey}
          />
        </div>
        <div className='decorator-row'>
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
          { data.points.color.type !== 'constant' &&

            <FormSelectItem
              labelText={'Based on'}
              name={'color-data-key'}
              value={data.points.color.data_key || ''}
              inline={true}
              firstOptionText={'Select'}
              options={target_attrs}
              errors={errors}
              onChange={this.updateColorDataKey}
            />

          }
          <Dropdown className='color-dropdown' ref='color-dropdown'>
            <DropdownTrigger className='trigger'>
              { colorDropdownTrigger }
            </DropdownTrigger>
            <DropdownContent className='dropdown-contents'>
              { colorDropdownContent }
            </DropdownContent>
          </Dropdown>
        </div>
        <div className='decorator-row'>
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
          { data.points.size.type === 'constant' &&
            <FormItem
              labelText={'Value'}
              name={'value'}
              value={data.points.size.bounds[0]}
              inline={true}
              errors={errors}
              onChange={e => this.setSize(e)}
            />
          }
          { data.points.size.type !== 'constant' &&
            <div>
              <FormSelectItem
                labelText={'Based on'}
                name={'side-data-key'}
                value={data.points.size.data_key || ''}
                inline={true}
                firstOptionText={'Select'}
                options={target_attrs}
                errors={errors}
                onChange={this.updateSizeDataKey}
              />
              <FormItem
                labelText={'Value'}
                name={'value'}
                value={data.points.size.bounds[0]}
                inline={true}
                errors={errors}
                onChange={e => this.setLowerSize(e)}
              />
              <FormItem
                labelText={'Max'}
                name={'upper-size'}
                value={data.points.size.bounds[1]}
                inline={true}
                errors={errors}
                onChange={e => this.setUpperSize(e)}
              />
            </div>
          }
        </div>
        <div className='decorator-row'>
          <FormItem
            labelText={'Sprite'}
            name={'sprite'}
            value={data.points.sprite}
            errors={errors}
            onChange={this.setSprite}
          />
        </div>
        <button onClick={this.save}>Save Viz</button>
      </div>
    )
  }
}
