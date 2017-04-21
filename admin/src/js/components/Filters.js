/* @flow */

import React, { Component } from 'react'
import Dropdown, { DropdownTrigger, DropdownContent } from 'react-simple-dropdown';

import { HamburgerIcon, OpenInNewIcon, ArrowDownIcon } from './icons/Icons'

import type {StringFilter, DateFilter, NumFilter} from './Collection'

export class StringFilterComponent extends Component {
    props: {
        data: StringFilter;
        creator: Object;
    }

    constructor(props: StringFilter) {
        super(props);
    }

    render() {
        const operations = ["contains","does not contain","matches","exists"].map((o,i) => <option value={o} key={i}>{o.toUpperCase()}</option>)
        const {data,creator} = this.props
        let value_field 
        if(data.options.length > 0){
            let options = data.options.map((name,i) => {
                return (
                    <option value={name} key={i}>{name}</option>
                )
            })
            value_field = (
                <select className="value-body-select" value={data.query} onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}>
                    {options}
                </select>
            )
        } else {
            value_field = <input className="filter-body-input" value={data.query} onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}/>
        }
        
        return (
            <div className="fk-filter fk-guidfilter" key={data.id}>
                <div className="filter-title-bar">
                    <span>{data.attribute}</span>
                    <div className="filter-title-controls">
                        <span className="filter-icon"></span>
                        <span className="filter-closer"></span>
                    </div>
                </div>
                <div className="filter-body">
                    <div>
                        <span className="filter-body-label">Operation: </span>
                        <select className="filter-body-select" value={data.operation} onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}>
                            {operations}
                        </select>
                    </div>
                    <div>
                        <span className="filter-body-label">Value: </span>
                        {value_field}
                    </div>
                    <div className="filter-body-buttons">
                        <button className="filter-body-cancel" onClick={() => creator.deleteFilter(data)}>Delete</button>
                        <button className="filter-body-save">Save</button>
                    </div>
                </div>
            </div>
        )
    }
}

export class NumFilterComponent extends Component {
    props: {
        data: NumFilter;
        creator: Object;
    }

    constructor(props: NumFilter) {
        super(props);
    }

    render() {
        const operations = ["GT","LT","EQ","notch"].map((o,i) => <option value={o} key={i}>{o.toUpperCase()}</option>)
        const {creator,data} = this.props
        
        return (
            <div className="fk-filter fk-guidfilter">
                <div className="filter-title-bar">
                    <span>{data.attribute}</span>
                    <div className="filter-title-controls">
                        <span className="filter-icon"></span>
                        <span className="filter-closer"></span>
                    </div>
                </div>
                <div className="filter-body">
                    <div>
                        <span className="filter-body-label">Operation: </span>
                        <select className="filter-body-select" value={data.operation} onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}>
                            {operations}
                        </select>
                    </div>
                    <div>
                        <span className="filter-body-label">Value: </span>
                        <input className="filter-body-input" value={data.query} onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}/>
                    </div>
                    <div className="filter-body-buttons">
                        <button className="filter-body-cancel" onClick={() => creator.deleteFilter(data)}>Delete</button>
                        <button className="filter-body-save">Save</button>
                    </div>
                </div>
            </div>
        )
    }
}

export class DateFilterComponent extends Component {
    props: {
        data: DateFilter;
        creator: Object;
    }

    constructor(props: NumFilter) {
        super(props);
    }

    render() {
        const operations = ["before","after","within"].map((o,i) => <option value={o} key={i}>{o.toUpperCase()}</option>)
        const {data,creator} = this.props
        const withinput = data.operation === "within" ?
                            (
                                <div>
                                    <span className="filter-body-label">Within: </span>
                                    <input className="filter-body-input" value={data.within}/>
                                </div>
                            ) : null
        
        return (
            <div className="fk-filter fk-guidfilter">
                <div className="filter-title-bar">
                    <span>{data.attribute}</span>
                    <div className="filter-title-controls">
                        <span className="filter-icon"></span>
                        <span className="filter-closer"></span>
                    </div>
                </div>
                <div className="filter-body">
                    <div>
                        <span className="filter-body-label">Operation: </span>
                        <select className="filter-body-select" value={data.operation} onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}>
                            {operations}
                        </select>
                    </div>
                    <div>
                        <span className="filter-body-label">Date: </span>
                        <input className="filter-body-input" value={data.date} onChange={(e) => creator.updateFilter(data,{"date":e.target.value})}/>
                    </div>
                    {withinput}
                    <div className="filter-body-buttons">
                        <button className="filter-body-cancel" onClick={() => creator.deleteFilter(data)}>Delete</button>
                        <button className="filter-body-save">Save</button>
                    </div>
                </div>
            </div>
        )
    }
}
