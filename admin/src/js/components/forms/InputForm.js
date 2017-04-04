// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APIUser, APIFieldkitInput, APINewFieldkitInput, APIMutableInput, APINewTwitterInput, APITeam } from '../../api/types';

type Props = {
  input?: ?APIMutableInput,
  users: {[id: number]: APIUser},
  teams: APITeam[],

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (i: APIMutableInput) => Promise<?APIErrors>;
}

export class InputForm extends Component {
  props: Props;

  /*flow-include
  state: {
    ...$Exact<APIMutableInput>,  
    users: {[id: number]: APIUser},
    teams: APITeam[],
    bindingType: ?string,
    saveDisabled: boolean,
    errors: ?APIErrors
  };
  */  
  constructor(props: Props) {
    super(props)

    this.state = {
      team_id: -1,
      user_id: -1,
      name: '',
      ...props.input,
      users: props.users,
      teams: props.teams,
      bindingType:  this.getBindingType(props),
      saveDisabled: false,
      errors: null
    }
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setState({
      team_id: -1,
      user_id: -1,
      name: '',
      ...nextProps.input,
      users: nextProps.users,
      teams: nextProps.teams,
      bindingType: this.getBindingType(nextProps),
      saveDisabled: false,
      errors: null
    });
  }

  resetInputBinding(bindingType: string) {
    if(bindingType === 'team') {
      this.setState({user_id: undefined});
    } else if (bindingType === 'member') {
      this.setState({team_id: undefined});
    }
  }

  getBindingType(props: Props): ?string{
    const { input } = this.props;
    if(input){
      if(input.team_id) {
        return 'team';
      } else if (input.user_id) {
        return 'member';
      }
    }
  }

  async save(inputId?: number) {
    const { name, team_id, user_id } = this.state;
    const errors = await this.props.onSave({
      name: name,
      team_id: parseInt(team_id),
      user_id: parseInt(user_id)
    }, inputId);
    if (errors) {
      this.setState({ errors });
    }
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({ [name]: value });

    if(name === "bindingType"){
      this.resetInputBinding(value);
    }    
  }

  render () {
    const { input } = this.props;
    const { users, teams, name, bindingType, team_id, user_id, errors } = this.state;

    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input type="text" name="name" className="lg" value={name} onChange={this.handleInputChange.bind(this)} />
          { errorsFor(errors, 'name') }
        </div>

        { input &&
          <div>
            <div className="form-group">
              <label htmlFor="bindingType">Bind input to</label>
              <select name="bindingType"  value={bindingType} onChange={this.handleInputChange.bind(this)}>
                <option disabled selected>Select</option>
                <option value={"team"}>Team</option>
                <option value={"member"}>Member</option>
              </select>
            </div>

            { bindingType === "team" &&
              <div className="form-group">
                <label htmlFor="team_id">Team</label>
                <select name="team_id" value={team_id} onChange={this.handleInputChange.bind(this)}>
                  <option disabled selected>Select a team</option>
                  { teams.map((team, i) => 
                    <option key={i} value={team.id}>{team.name}</option>) }
                </select>
                { errorsFor(errors, 'team_id') }          
              </div>
            }
            { bindingType === "member" &&
              <div className="form-group">
                <label htmlFor="user_id">Member</label>
                <select name="user_id" value={user_id} onChange={this.handleInputChange.bind(this)}>
                  <option disabled selected>Select a member</option>
                  { Object.keys(users).map(id => 
                    <option key={id} value={id}>{users[id].username}</option>) }
                </select>
                { errorsFor(errors, 'user_id') }          
              </div>
            }
          
          </div>
        }

      </FormContainer>
    )
  }
}