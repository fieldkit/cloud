// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { MembersTable } from '../shared/MembersTable';
import { TeamForm } from '../forms/TeamForm';
import { MemberForm } from '../forms/MemberForm';
import { FKApiClient } from '../../api/api';
import '../../../css/teams.css'

import { FormContainer } from '../containers/FormContainer';

import type { APIProject, APIExpedition, APITeam, APINewTeam, APINewMember, APIMember, APIUser } from '../../api/types';

type Props = {
  project: APIProject;
  expedition: APIExpedition;

  match: Object;
  location: Object;
  history: Object;
}

export class Teams extends Component {

  props: Props;
  state: {
    teams: APITeam[],
    members: { [teamId: number]: APIMember[] },
    users: { [teamId: number]: APIUser[] },
    teamDeletion: ?{
      contents: React$Element<*>;
      teamId: number;
    },
    memberDeletion: ?{
      contents: React$Element<*>;
      userId: number;
      teamId: number;
    }
  }

  constructor (props: Props) {
    super(props);
    this.state = {
      teams: [],
      members: {},
      users: {},
      teamDeletion: null,
      memberDeletion: null
    }

    this.loadTeams();
  }

  async loadTeams() {
    const teamsRes = await FKApiClient.get().getTeamsBySlugs(this.props.project.slug, this.props.expedition.slug);
    if (teamsRes.type === 'ok' && teamsRes.payload) {
      const { teams } = teamsRes.payload;
      this.setState({ teams: teams });
      for (const team of teams) {
        await this.loadMembers(team.id);
      }
    }
  }

  async loadMembers(teamId: number) {
    const membersRes = await FKApiClient.get().getMembers(teamId);
    if (membersRes.type === 'ok' && membersRes.payload) {
      const { members } = this.state;
      members[teamId] = membersRes.payload.members;
      this.setState({members: members});
      for (const member of members[teamId]) {
        await this.loadMemberName(teamId, member.user_id);
      }      
    }
  }

  async loadMemberName(teamId: number, userId: number){
    const userRes = await FKApiClient.get().getUserById(userId);
    if (userRes.type === 'ok' && userRes.payload) {
      const { users } = this.state;
      if(!users[teamId]){
        users[teamId] = [];
      }users[teamId].push(userRes.payload);
      this.setState({users: users});
    }    
  }  

  async onTeamCreate(e: APINewTeam) {

    const { expedition, match } = this.props;

    const teamRes = await FKApiClient.get().createTeam(expedition.id, e);
    if (teamRes.type === 'ok') {
      await this.loadTeams();
      this.props.history.push(`${match.url}`);
    } else {
      return teamRes.errors;
    }
  }

  startTeamDelete(e: APITeam) {
    const teamId = e.id;
    this.setState({
      teamDeletion: {
        contents: <span>Are you sure you want to delete the <strong>{e.name}</strong> team?</span>,
        teamId
      }
    })
  }

  async confirmTeamDelete() {
    const { teamDeletion } = this.state;

    if (teamDeletion) {
      const { teamId } = teamDeletion;

      const teamRes = await FKApiClient.get().deleteTeam(teamId);
      if (teamRes.type === 'ok') {
        await this.loadTeams(teamId);
        this.setState({ teamDeletion: null })
      } else {
        // TODO: show errors somewhere
      }
    }
  }
  
  cancelTeamDelete() {
    this.setState({ teamDeletion: null });
  }

  async onMemberAdd(teamId: number, e: APINewMember) {
    const { match } = this.props;
    const memberRes = await FKApiClient.get().addMember(teamId, e);
    if (memberRes.type === 'ok') {
      await this.loadTeams();
      this.props.history.push(`${match.url}`);
    } else {
      return memberRes.errors;
    }
  }

  startMemberDelete(teamId: number, userId: number) {
    const { teams, members, users } = this.state;
    let teamName = 'unknownTeam';
    let userName = 'unknownUser';

    if (members[teamId]) {
      const teamMembers = members[teamId];
      const teamMember = teamMembers.find((u: APIMember) => u.user_id === userId);
      const team = teams.find((t: APITeam) => t.id === teamId);

      if (team) {
        teamName = team.name;
      }
      if (teamMember) {
        const user = users[teamId].find((u: APIUser) => u.id === teamMember.user_id);
        if (user) {
          userName = user.name;
        }
      }
    }
    this.setState({
      memberDeletion: {
        contents: <span>Are you sure you want to remove <strong>{userName}</strong> from <strong>{teamName}?</strong></span>,
        userId, teamId
      }
    })
  }

  async confirmMemberDelete() {
    const { memberDeletion } = this.state;

    if (memberDeletion) {
      const { teamId, userId } = memberDeletion;

      const memberRes = await FKApiClient.get().deleteMember(teamId, userId);
      if (memberRes.type === 'ok') {
        await this.loadMembers(teamId);
        this.setState({ memberDeletion: null })
      } else {
        // TODO: show errors somewhere
      }
    }
  }

  cancelMemberDelete() {
    this.setState({ memberDeletion: null });
  }

  startMemberUpdate(teamId: number, memberId: number, role: string) {
    const { members } = this.state;
    const member = members[teamId].find(member => member.user_id === memberId);
    if(member) {
      member.role = role;      
    }
    this.setState({members: members});
  }

  async onTeamUpdate(teamId: number, team: APINewTeam) {
    const teamRes = await FKApiClient.get().updateTeam(teamId, team);
    if(teamRes.type === 'ok' && teamRes.payload) {
      console.log(teamRes.payload);
    } else {
      return teamRes.errors;
    }
  }

  getTeamById (id: number): ?APITeam {
    const { teams } = this.state;
    return this.state.teams.find(team => team.id === id);
  }

  render() {
    const { match } = this.props;
    const { teams, members, users, memberDeletion, teamDeletion } = this.state;

    return (
      <div className="teams">

        <Route path={`${match.url}/new-team`} render={() =>
          <ReactModal isOpen={true} contentLabel="Create New Team" className="modal" overlayClassName="modal-overlay">
            <h2>Create New Team</h2>
            <TeamForm
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onTeamCreate.bind(this)} />
          </ReactModal> } />

        <Route path={`${match.url}/:teamId/add-member`} render={props =>
          <ReactModal isOpen={true} contentLabel="Add Members" className="modal" overlayClassName="modal-overlay">
            <h2>Add Member</h2>
            <MemberForm
              teamId={props.match.params.teamId}
              members={this.state.members[props.match.params.teamId]}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onMemberAdd.bind(this)} 
              saveText="Add" />
          </ReactModal> } />          

        <Route path={`${match.url}/:teamId/edit`} render={props =>
          <ReactModal isOpen={true} contentLabel="Edit Team" className="modal" overlayClassName="modal-overlay">
            <TeamForm
              team={this.getTeamById(parseInt(props.match.params.teamId))}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onTeamUpdate.bind(this, parseInt(props.match.params.teamId))} />
          </ReactModal> } />          

        { teamDeletion &&
          <ReactModal isOpen={true} contentLabel="Delete Team" className="modal" overlayClassName="modal-overlay">
            <h2>Delete Team</h2>
            <FormContainer
              onSave={this.confirmTeamDelete.bind(this)}
              onCancel={this.cancelTeamDelete.bind(this)}
              saveText="Confirm"
            >
              <div>{teamDeletion.contents}</div>
            </FormContainer>
          </ReactModal> }

        { memberDeletion &&
          <ReactModal isOpen={true} contentLabel="Remove Member" className="modal" overlayClassName="modal-overlay">
            <h2>Remove Member</h2>
            <FormContainer
              onSave={this.confirmMemberDelete.bind(this)}
              onCancel={this.cancelMemberDelete.bind(this)}
              saveText="Confirm"
            >
              <div>{memberDeletion.contents}</div>
            </FormContainer>
          </ReactModal> }

        <h1>Teams</h1>

        <div className="accordion">
        { teams.map((team, i) =>
          <div key={i} className="accordion-row expanded">
              <div className="accordion-row-header">
                <h4>{team.name}</h4>
                <div className="nav">
                  <Link className="button secondary" to={`${match.url}/${team.id}/edit`}>Edit</Link>
                  <button className="secondary" onClick={this.startTeamDelete.bind(this, team)}>Delete</button>
                </div>
              </div>
              <p>{team.description}</p>
              { members[team.id] && members[team.id].length > 0 &&
                <MembersTable
                  teamId={team.id}
                  members={members[team.id]}
                  users={users[team.id]}
                  onDelete={this.startMemberDelete.bind(this)} 
                  onMemberUpdate={this.startMemberUpdate.bind(this)}/> }
              { (!members[team.id] || members[team.id].length === 0) &&
                <p className="empty">This team has no members yet.</p> }
              <Link className="button secondary" to={`${match.url}/${team.id}/add-member`}>Add Member</Link>                

          </div> ) }
        { teams.length === 0 &&
          <span className="empty">No teams</span> }
        </div>
        <Link className="button" to={`${match.url}/new-team`}>Create New Team</Link>
      </div>
    )
  }
}
