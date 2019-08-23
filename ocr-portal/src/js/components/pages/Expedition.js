// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal'

import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
import { FKApiClient } from '../../api/api';

import { RemoveIcon } from '../icons/Icons'

import type { APIProject, APIExpedition, APINewExpedition, APIInputToken, APINewInputToken } from '../../api/types';

type Props = {
project: APIProject;
expedition: APIExpedition;
onUpdate: (newSlug: ?string) => void;

match: Object;
location: Object;
history: Object;
}

export class Expedition extends Component {
    props: Props;
    state: {
    inputTokens: APIInputToken[];
    }

    constructor(props: Props) {
        super(props);
        this.state = {
            inputTokens: []
        }
        this.loadInputTokens();
    }

    async loadInputTokens() {
        const tokenRes = await FKApiClient.get().getExpeditionInputTokens(this.props.expedition.id);
        if (tokenRes.type === 'ok') {
            this.setState({
                inputTokens: tokenRes.payload.input_tokens
            });
        }
    }

    async onExpeditionSave(expedition: APINewExpedition) {
        const expeditionRes = await FKApiClient.get().updateExpedition(this.props.expedition.id, expedition);
        if (expeditionRes.type !== 'ok') {
            return expeditionRes.errors;
        }

        if (expeditionRes.slug != this.props.expedition.slug && expeditionRes.payload) {
            this.props.onUpdate(expeditionRes.payload.slug);
        } else {
            this.props.onUpdate();
        }
    }

    async deleteInputToken(inputTokenId: number) {
        const res = await FKApiClient.get().deleteInputToken(inputTokenId);
        if (res.type !== 'ok') {
            // TODO: show error
        }

        await this.loadInputTokens();
    }

    async addInputToken() {
        const {expedition} = this.props;
        const res = await FKApiClient.get().createInputToken(expedition.id, {});
        if (res !== 'ok') {
            // TODO: show error
        }

        await this.loadInputTokens();
    }

    render() {
        const {project, expedition} = this.props;
        const {inputTokens} = this.state;

        return (
            <div className="expedition">
                <h1>Expedition Settings</h1>
                <div className="row">
                    <div className="two-columns">
                        <ProjectExpeditionForm projectSlug={ project.slug } name={ expedition.name } slug={ expedition.slug } description={ expedition.description } onSave={ this.onExpeditionSave.bind(this) }
                        />
                    </div>
                </div>
                <div className="row">
                    <h3>Input Tokens</h3>
                    { inputTokens.length > 0 &&
                      <table className="input_tokens-table">
                          <thead>
                              <tr>
                                  <th>Token</th>
                                  <th></th>
                              </tr>
                          </thead>
                          <tbody>
                              { inputTokens.map((token, i) => <tr key={ i } className="input_token-item">
                                                                  <td>
                                                                      { token.token }
                                                                  </td>
                                                                  <td>
                                                                      <div className="bt-icon medium" onClick={ () => this.deleteInputToken(token.id) }>
                                                                          <RemoveIcon />
                                                                      </div>
                                                                  </td>
                                                              </tr>) }
                          </tbody>
                      </table> }
                    { inputTokens.length === 0 && <div>No tokens!</div> }
                    <button className="button" onClick={ this.addInputToken.bind(this) }>Create Input Token</button>
                </div>
            </div>
        )
    }
}
