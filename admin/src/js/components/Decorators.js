/* @flow */
import React, { Component } from 'react'
import { SketchPicker } from 'react-color';
import ColorBrewer from 'colorbrewer';

export type Stop = {
    location: number;
    color: string;
}

function cloneStop(s: Stop) : Stop {
    const {location, color} = s;
    return {
        location,
        color
    }
}

export type Color = {
    type: "constant" | "linear";
    colors: Stop[];
    data_key: ?string;
    bounds: ?[number,number];
}

function cloneColor(c: Color): Color {
    const {type, data_key} = c;
    const colors = c.colors.map(s => cloneStop(s))
    const bounds = c.bounds ? [c.bounds[0],c.bounds[1]] : null
    return {
        type,
        colors,
        data_key,
        bounds
    }
}

export type Size = {
    type: "constant" | "linear";
    data_key: ?string;
    bounds: [number,number];
}

function cloneSize(s: Size): Size {
    const {type, data_key} = s;
    const bounds = [s.bounds[0],s.bounds[1]]

    return {
        type,
        data_key,
        bounds
    }
}

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

function clonePointDecorator(p: PointDecorator): PointDecorator{
    const {collection_id, title, type} = p
    const points = {
        color: cloneColor(p.points.color),
        size: cloneSize(p.points.size),
        sprite: p.points.sprite
    }
    return {
        collection_id,
        points,
        title,
        type
    }
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
    props: {initial_state: PointDecorator}
    state: {data: PointDecorator}

    constructor(props: PointDecoratorProps){
        super(props)
    }

    getInitialState(){
        return {data: this.props.initial_state}
    }

    setConstantColor(color: string){
        this.setState({})
    }

    render(){
       const {data} = this.state;
       const collections = null
       let color 
       if( data.points.color.type === "constant"){
         color = <SketchPicker onChangeComplete={c => this.setConstantColor(c)} color={data.points.color.colors[0]}/>
       } else {
       
       }
       const size = null
        
       return (
           <div className="point-decorator">
                <div className="decorator-row">
                    <span className="decorator-row-label">Source: </span>
                    <select>
                        {collections}
                    </select>
                </div>
                <div className="decorator-row">
                    <span className="decorator-row-label">Color: </span>
                    <select value={data.points.color.type}>
                        <option value="constant">constant</option>
                        <option value="linear">linear</option>
                    </select>
                    {color}
                </div>
                <div className="decorator-row">
                    <span className="decorator-row-label">Size: </span>
                    <select value={data.points.size.type}>
                        <option value="constant">constant</option>
                        <option value="linear">linear</option>
                    </select>
                </div>
                <div className="decorator-row">
                    <span className="decorator-row-label">Sprite: </span>
                    <input value={data.points.sprite}/>
                </div>
           </div>
       )
    }
}
