/* @flow */

import React, { Component } from 'react';
import { FormContainer } from '../containers/FormContainer';
import type { APIErrors, APICollection, APINewCollection, } from '../../api/types';
import { AddIcon } from '../icons/Icons';
import log from 'loglevel';
import { FKApiClient } from '../../api/api';
import { StringFilterComponent, NumFilterComponent, DateFilterComponent, } from '../forms/FiltersForm';
import type { Collection, Attr, Target, ProjectData, Filter, FilterFn, StringAttr, StringFilter, NumAttr, NumFilter, DateAttr, DateFilter, } from '../../types/CollectionTypes';
import type { APIProject } from '../../api/types';

function stringifyOptions(attr: Attr): string {
    if (attr.type === 'string') {
        return attr.options.join(', ');
    } else if (attr.type === 'num') {
        return attr.options.map(n => n.toString()).join(', ');
    } else if (attr.type === 'date') {
        return attr.options.map(n => n.toString()).join(', ');
    }

    return '';
}

function cloneCollection(c: Collection): Collection {
    return {
        name: c.name,
        id: c.id,
        filters: c.filters.slice(0),
        guid_filters: c.guid_filters.slice(0),
        string_filters: c.string_filters.slice(0),
        num_filters: c.num_filters.slice(0),
        geo_filters: c.geo_filters.slice(0),
        date_filters: c.date_filters.slice(0),
        mods: c.mods.slice(0),
        string_mods: c.string_mods.slice(0),
        num_mods: c.num_mods.slice(0),
        geo_mods: c.geo_mods.slice(0),
    };
}

function getNextId(c: Collection): number {
    const max = Math.max(...c.filters);
    return max > -1 ? max + 1 : 0;
}

function emptyStringFilter(c: Collection, attr: StringAttr): StringFilter {
    let next_id = getNextId(c);
    return {
        id: next_id,
        attribute: attr.name,
        operation: 'contains',
        query: '',
        options: attr.options,
        type: 'string',
        target: attr.target,
    };
}

function emptyNumFilter(c: Collection, attr: NumAttr): NumFilter {
    let next_id = getNextId(c);
    return {
        id: next_id,
        attribute: attr.name,
        operation: 'GT',
        query: 3,
        type: 'num',
        target: attr.target,
    };
}

function emptyDateFilter(c: Collection, attr: DateAttr): DateFilter {
    let next_id = getNextId(c);
    let f = {
        id: next_id,
        attribute: attr.name,
        operation: 'before',
        date: 0,
        within: 0,
        type: 'date',
        target: attr.target,
    };
    return f;
}

function emptyCollection(): Collection {
    return {
        name: '',
        id: '',
        filters: [],
        guid_filters: [],
        string_filters: [],
        num_filters: [],
        geo_filters: [],
        date_filters: [],
        mods: [],
        string_mods: [],
        num_mods: [],
        geo_mods: [],
    };
}

type Props = {
collection?: ?APICollection,

cancelText?: string,
saveText?: ?string,
onCancel?: () => void,
onSave: (c: APINewCollection, collectionId?: number) => Promise<?APIErrors>,
};

export class ProjectCollectionsForm extends Component {
    projectData: ProjectData;
    props: Props;
    state: {
    collection: Collection,
    errors: ?APIErrors,
    };

    constructor(props: Props) {
        super(props);

        this.projectData = {
            expeditions: {
                name: 'expedition',
                options: ['ITO 2015', 'ITO 2016', 'ITO 2017'],
                type: 'string',
                target: 'expedition',
            },
            bindings: [],
            doctypes: [],
            attributes: [
                {
                    name: 'message',
                    options: [],
                    type: 'string',
                    target: 'attribute',
                },
                {
                    name: 'user',
                    options: ['@eric', '@othererik', '@gabriel'],
                    type: 'string',
                    target: 'attribute',
                },
                {
                    name: 'humidity',
                    options: [],
                    type: 'num',
                    target: 'attribute',
                },
                {
                    name: 'created',
                    options: [],
                    type: 'date',
                    target: 'attribute',
                } ,
            ],
        } ;

        this.state = {
            collection: emptyCollection(),
            errors: null,
        } ;
    }

    componentWillReceiveProps(nextProps: Props) {}

    getAttrRow(
        attr: Attr,
        component_lookup: { [number]: Object },
        target: Target
    ): * {
        const attribute_row = (
        <div key={ 'filter-' + attr.name } className="accordion-row-header">
            <h4>{ attr.name }</h4>
            <button className="secondary" onClick={ () => this.addFilter(attr.name, target) }>
                <div className="bt-icon medium">
                    <AddIcon />
                </div>
                Add Filter
            </button>
        </div>
        );
        let buttons: Object[] = [attribute_row];
        this.state.collection.filters.forEach(f => {
            const component = component_lookup[f];
            if (component) {
                const {attribute} = component.props.data;
                if (attr.name === attribute) {
                    buttons.push(component);
                }
            }
        });
        return (
            <div className="accordion-row expanded" key={ attr.name }>
                { buttons }
            </div>
            );
    }

    addFilter(name: string, target: Target) {
        let attribute;
        let new_filter;
        const project_data = this.projectData;
        if (target === 'expedition') {
            attribute = project_data.expeditions;
        } else if (target === 'binding') {
            attribute = project_data.bindings.find(b => b.name === name);
        } else if (target === 'doctype') {
            attribute = project_data.doctypes.find(d => d.name === name);
        } else if (target === 'attribute') {
            attribute = project_data.attributes.find(a => a.name === name);
        } else {
            throw new Error('Incompatible target');
        }
        if (!attribute) {
            return;
        }
        let collection = cloneCollection(this.state.collection);
        if (attribute.type === 'string') {
            new_filter = emptyStringFilter(collection, attribute);
            collection.string_filters.push(new_filter);
        } else if (attribute.type === 'num') {
            new_filter = emptyNumFilter(collection, attribute);
            collection.num_filters.push(new_filter);
        } else if (attribute.type === 'date') {
            new_filter = emptyDateFilter(collection, attribute);
            collection.date_filters.push(new_filter);
        }
        if (new_filter) {
            collection.filters.push(new_filter.id);
            this.setState({
                collection
            });
        }
    }

    updateFilter = (filter, update) => {
        let new_filter = Object.assign(filter, update);
        let collection = cloneCollection(this.state.collection);
        if (new_filter.type === 'string') {
            collection.string_filters = collection.string_filters.filter(
                f => f.id !== new_filter.id
            );
            collection.string_filters.push(new_filter);
        } else if (filter.type === 'date') {
            collection.date_filters = collection.date_filters.filter(
                f => f.id !== new_filter.id
            );
            collection.date_filters.push(new_filter);
        } else if (new_filter.type === 'num') {
            collection.num_filters = collection.num_filters.filter(
                f => f.id !== new_filter.id
            );
            collection.num_filters.push(new_filter);
        }
        this.setState({
            collection
        });
    };

    deleteFilter(f: Filter) {
        let collection = cloneCollection(this.state.collection);
        if (f.type === 'string') {
            collection.string_filters = collection.string_filters.filter(
                cf => cf.id !== f.id
            );
        } else if (f.type === 'num') {
            collection.num_filters = collection.num_filters.filter(
                cf => cf.id !== f.id
            );
        } else if (f.type === 'date') {
            collection.date_filters = collection.date_filters.filter(
                cf => cf.id !== f.id
            );
        }
        collection.filters = collection.filters.filter(cf => cf !== f.id);
        this.setState({
            collection
        });
    }

    save() {
        const {collection} = this.state;
        console.log(JSON.stringify(collection, null, ' '));
    // TO-DO
    // const errors = await this.props.onSave(collection);
    // if (errors) {
    //   this.setState({ errors });
    // }
    }

    render() {
        const projectData = this.projectData;
        const {errors} = this.state;
        let {collection} = this.state;

        let component_lookup: { [number]: Object } = {};

        collection.num_filters.forEach((f, i) => {
            component_lookup[f.id] = (
                <NumFilterComponent creator={ this } data={ f } key={ 'n-' + i } errors={ errors } />
            );
        });
        collection.string_filters.forEach((f, i) => {
            component_lookup[f.id] = (
                <StringFilterComponent creator={ this } data={ f } key={ 's-' + i } errors={ errors } />
            );
        });
        collection.date_filters.forEach((f, i) => {
            component_lookup[f.id] = (
                <DateFilterComponent creator={ this } data={ f } key={ 'd-' + i } errors={ errors } />
            );
        });

        let filters_by_attribute = projectData.attributes.map(attr => {
            return this.getAttrRow(attr, component_lookup, 'attribute');
        }, {});

        let filters_by_binding = projectData.bindings.map(attr => {
            return this.getAttrRow(attr, component_lookup, 'binding');
        }, {});

        let filters_by_doctype = projectData.doctypes.map(attr => {
            return this.getAttrRow(attr, component_lookup, 'doctype');
        }, {});

        let expedition_filter = this.getAttrRow(
            projectData.expeditions,
            component_lookup,
            'expedition'
        );

        return (
            <FormContainer onSave={ this.save.bind(this) } onCancel={ this.props.onCancel } saveText={ this.props.saveText } cancelText={ this.props.cancelText }>
                <div className="accordion">
                    { expedition_filter }
                </div>
                <div className="accordion">
                    { filters_by_binding }
                </div>
                <div className="accordion">
                    { filters_by_doctype }
                </div>
                <div className="accordion">
                    { filters_by_attribute }
                </div>
            </FormContainer>
            );
    }
}
