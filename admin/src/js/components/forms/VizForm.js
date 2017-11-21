// @flow weak
import React, { Component } from 'react';
import ReactModal from 'react-modal';
import type { Lens_ } from 'safety-lens';
import { set } from 'safety-lens'
import { prop } from 'safety-lens/es2015'
import VizDecoratorPointForm from './VizDecoratorPointForm';
import VizGroupingForm from './VizGroupingForm';
import VizSelectionForm from './VizSelectionForm';
import type { ProjectData } from '../../types/CollectionTypes';
import type { Viz, SelectionOperation } from '../../types/VizTypes';
import type { APIErrors } from '../../api/types';

export const _selectionOperations: Lens_<Viz,SelectionOperation[]> = prop('selection_operations');

function updateViz<A>(l: Lens_<Viz, A>, value: A, viz: Viz): Viz {
    return set(l, value, viz);
}

type VizProps = {
initial_state: Viz,
project_data: ProjectData,
};

export default class VizComponent extends Component {
    props: VizProps;
    state: {
    data: Viz,
    errors: ?APIErrors,
    modal_open: boolean,
    };

    constructor(props: VizProps) {
        super(props);
        this.state = {
            data: this.props.initial_state,
            errors: null,
            modal_open: false,
        } ;
    }

    update<A>(lens: Lens_<Viz, A>, value: A): void {
        let {data} = this.state;
        data = updateViz(lens, value, data);
        this.setState({
            data
        });
    }

    getCollectionAttributes(): Attr[] {
        // TODO: CONNECT TO ACTUAL COLLECTIONS
        return this.props.project_data.attributes;
    }

    newSelectionID(): number {
        const {data} = this.state;
        const selections = data.selection_operations.map(s => s.id);
        const max = Math.max(...selections);
        return max > -1 ? max + 1 : 0;
    }

    addSelection(selection: SelectionOperation) {
        let selections = this.state.data.selection_operations.slice(0);
        selections.push(selection);
        this.update(_selectionOperations, selections);
        this.setState({
            modal_open: false
        });
    }

    deleteSelection(selection_id: number) {
        let selections = this.state.data.selection_operations.slice(0);
        selections = selections.filter(s => s.id !== selection_id);
        this.update(_selectionOperations, selections);
    }

    render() {
        const {data, errors, modal_open} = this.state;
        const default_attribute = this.getCollectionAttributes()[0];
        const new_selection = {
            id: this.newSelectionID(),
            value_name: '',
            source_attribute: default_attribute,
            operation: 'avg',
        };
        const selections = data.selection_operations.map((s, i) => {
            return (
                <span key={ i } className='selection' data-selection-id={ s.id }>
                                          { s.value_name }
                                          <span
                                className='selection-deleter'
                                onClick={ () => this.deleteSelection(s.id) }
                                >
                                            ×
                                          </span>
                </span>
                );
        });
        let decorator_component;
        if (data.decorator.type === 'point') {
            decorator_component = (
                <VizDecoratorPointForm viz={ data } creator={ this } />
            );
        }

        return (
            <div>
                <div>
                    <VizGroupingForm data={ data } errors={ errors } creator={ this } />
                </div>
                <div>
                    <span>Selections:</span>
                    <div>
                        { selections }
                    </div>
                    <button onClick={ () => this.setState({
                                          modal_open: !modal_open
                                      }) }>
                        Add Selection
                    </button>
                    <ReactModal contentLabel="New Selection" isOpen={ modal_open }>
                        <VizSelectionForm data={ data } initial_state={ new_selection } errors={ errors } creator={ this } />
                    </ReactModal>
                </div>
                { decorator_component }
            </div>
            );
    }
}
