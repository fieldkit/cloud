// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { FormItem } from './FormItem';
import { errorsFor, slugify, fkProtocol, fkHost } from '../../common/util';

import type { APIErrors, APINewExpedition } from '../../api/types';

type Props = {
projectSlug: string,
name?: string,
slug?: string,
description?: string,

cancelText?: string;
saveText?: ?string;
onCancel?: () => void;
onSave: (e: APINewExpedition) => Promise<?APIErrors>;
}

export class ProjectExpeditionForm extends Component {
    props: Props;
    state: {
    name: string,
    slug: string,
    description: string,
    slugHasChanged: boolean,
    saveDisabled: boolean,
    errors: ?APIErrors
    };

    constructor(props: Props) {
        super(props)
        const {name, slug, description} = props;

        this.state = {
            name: name || '',
            slug: slug || '',
            description: description || '',
            slugHasChanged: !!name && !!slug && slugify(name) != slug,
            saveDisabled: true,
            errors: null
        }
    }

    componentWillReceiveProps(nextProps: Props) {
        const {name, slug, description} = nextProps;

        this.setState({
            name: name || '',
            slug: slug || '',
            description: description || '',
            slugHasChanged: !!name && !!slug && slugify(name) != slug,
            saveDisabled: true,
            errors: null
        });
    }

    async save() {
        const errors = await this.props.onSave({
            name: this.state.name,
            slug: this.state.slug,
            description: this.state.description
        });
        if (errors) {
            this.setState({
                errors
            });
        }
    }

    handleNameChange(event) {
        const v = event.target.value;
        let stateUpdate = {
            name: v,
            saveDisabled: false
        };
        if (!this.state.slugHasChanged) {
            stateUpdate = {
                slug: slugify(v),
                ...stateUpdate
            } ;
        }
        this.setState(stateUpdate);
    }

    handleSlugChange(event) {
        const v = event.target.value;
        this.setState({
            slug: v,
            slugHasChanged: true,
            saveDisabled: false
        });
    }

    handleInputChange(event) {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        this.setState({
            [name]: value,
            saveDisabled: false
        });
    }

    render() {
        const {projectSlug} = this.props;
        const {name, slug, description, errors, saveDisabled} = this.state;

        return (
            <FormContainer onSave={ this.save.bind(this) } onCancel={ this.props.onCancel } saveText={ this.props.saveText } cancelText={ this.props.cancelText } saveDisabled={ saveDisabled }>
                <FormItem name="name" value={ name } labelText="name" className="lg" onChange={ this.handleNameChange.bind(this) } errors={ errors } />
                <div className="url-preview">
                    <p className="label">Your expedition will be available at the following address:</p>
                    <p className="url">
                        { `${fkProtocol()}//${projectSlug}.${fkHost()}/` }
                        <input type="text" name="slug" className='slug' value={ slug } onChange={ this.handleSlugChange.bind(this) } />
                    </p>
                    { errorsFor(errors, 'slug') }
                </div>
                <FormItem name="description" value={ description } labelText="description" className="lg" onChange={ this.handleInputChange.bind(this) } errors={ errors } />
            </FormContainer>
        )
    }
}
