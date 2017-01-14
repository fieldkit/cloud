import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select'

class NewTeamsSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      addMemberValue: null,
      inputValues: {}
    }
  }

  render () {

    const { 
      currentProjectID,
      expedition,
      teams,
      members,
      currentTeam,
      editedTeam,
      setCurrentTeam,
      addTeam,
      startEditingTeam,
      stopEditingTeam,
      setTeamProperty,
      setMemberProperty,
      saveChangesToTeam,
      clearChangesToTeam,
      fetchSuggestedMembers,
      addMember,
      removeMember
    } = this.props

    const roleOptions = [
      'Expedition Leader', 'Team Leader', 'Team Member'
    ]
    const teamTabs = teams
      .map((t, i) => {
        let className = 'team-name '
        if (t === currentTeam) className += 'editable active '
        if (currentTeam.get('name') === '') className += 'required '
        console.log('ok', t.get('name'), currentTeam.get('name'))
        return (
          <li 
            className={ className }
            key={t.get('id')}
            onClick={() => { 
              if (t !== currentTeam) setCurrentTeam(t.get('id'))
            }}
          >
            <ContentEditable
              html={(() => {
                return this.props.teams.get(i).get('name')
                // if (!!t.get('name')) return t.get('name')
                // else if (t.get('status') === 'new') return 'New Team Name'
                // else return ''
              })()}
              disabled={t !== currentTeam}
              onClick={(e) => {
                if (t === currentTeam) {
                  if (this.props.teams.get(i).get('name') === 'New Team') setTeamProperty('name', '')
                  startEditingTeam()
                }
              }}
              onBlur={(e) => {
                if (this.props.teams.get(i).get('name') === '') setTeamProperty('name', 'New Team')
                stopEditingTeam()
              }}
              onChange={(e) => {
                setTeamProperty('name', e.target.value)
              }}
            />
          </li>
        )
      })

    const teamMembers = members
      .map((m, i) => {
        const inputs = m.get('inputs')
          .map((d, j) => {
            return <li className="tag" key={j}>{d}</li>
          })

        return (
          <tbody key={i} onClick={() => {
            // this.selectMember(i)
          }}>
            <td className="name">{m.get('name')}</td>
            <td className="role">
              <Dropdown
              options={roleOptions}
              onChange={(e) => {
                startEditingTeam()
                setMemberProperty(m.get('id'), 'role', e.value)
              }}
              value={currentTeam.getIn(['members', m.get('id'), 'role'])}
              placeholder="Select an option"
            />
            </td>
            <td className="inputs">
              <Select
                name="select-inputs"
                value={this.state.inputValues[m.get('id')]}
                multi={true}
                options={
                  [
                    { value: 'ambit', label: 'Ambit Tracker' },
                    { value: 'twitter', label: 'Twitter' },
                    { value: 'sightings', label: 'Sighting App' }
                  ]
                }
                onChange={(values) => {
                  this.setState({
                    ...this.state,
                    inputValues: {
                      ...this.state.inputValues,
                      [m.get('id')]: values
                    }
                  })
                }}
                clearable={false}
              />
            </td>
            <td className="activity">
              <svg></svg>
            </td>
            <td 
              className="remove"
              onClick={() => {
                startEditingTeam()
                removeMember(m.get('id'))
              }}
            >  
              <img src="/src/img/icon-remove-small.png"/>
            </td>
          </tbody>
        )
        // } else {
        //   return (
        //     <tbody key={i} onClick={() => {
        //       this.unselectMember()
        //     }}>
        //       <td className="name" colSpan="6" width="100%">
        //         <h2>{m.get('name')}</h2>
        //       </td>
        //     </tbody>
        //   )
        // }
      })

    const teamActionButtons = () => {
      const actionButtons = []

      if (currentTeam.get('new')) {
        actionButtons.push(
          <div
            className="button secondary"
            onClick={() => {
              this.props.removeCurrentTeam()
            }}
          >
            Cancel
          </div>
        )
        actionButtons.push(
          <div
            className="button secondary"
            onClick={() => {
              this.props.saveChangesToTeam()
            }}
          >
            Save New Team
          </div>
        )
      } else {
        actionButtons.push(
          <div
            className="button secondary"
            onClick={() => {
              this.props.removeCurrentTeam()
            }}
          >
            Remove Team
          </div>
        )

        if (!!currentTeam && !!editedTeam && 
            !!editedTeam.find((val, key) => {
              return currentTeam.get(key) !== val
            })) {
          actionButtons.push(
            <div
              className="button secondary"
              onClick={() => {
                this.props.saveChangesToTeam()
              }}
            >
              Save Changes
            </div>
          )
          if (editedTeam.get('status') !== 'new') {
            actionButtons.push(
              <div
                className="button secondary"
                onClick={() => {
                  this.props.clearChangesToTeam()
                }}
              >
                Clear Changes
              </div>
            )
          }
        }
      }
      return (
        <div className="actions">
          { actionButtons }
        </div>
      )
    }

    const selectedTeam = (
      <div className="team" key={currentTeam.get('id')}>
          { teamActionButtons() }
          <div className="header">
            <div className="column description">
              <h5>Description</h5>
              <ContentEditable
                html={(() => {
                  return currentTeam.get('description')
                })()}
                disabled={false}
                onClick={(e) => {
                  if (this.props.currentTeam.get('description') === 'Enter a description') setTeamProperty('description', '')
                  startEditingTeam()
                }}
                onBlur={(e) => {
                  if (this.props.currentTeam.get('description') === '') setTeamProperty('description', 'Enter a description')
                  stopEditingTeam()
                }}
                onChange={(e) => { 
                  setTeamProperty('description', e.target.value)
                }}
              />
            </div>
            <svg className="activity column">
            </svg>
          </div>
          <h5>Members</h5>
          <table className="objects-list">
            {
              !!members && !!members.size &&
              <tbody>
                <td className="name">Name</td>
                <td className="role">Role</td>
                <td className="inputs">Inputs</td>
                <td className="activity">Activity</td>
                <td className="remove"></td>
              </tbody>
            }
            { teamMembers }
            <tbody>
              <td className="add-object" colSpan="3" width="50%">
                <div className="add-object-container">
                  <Select.Async
                    name="add-member"
                    loadOptions={(input, callback) =>
                      fetchSuggestedMembers(input, callback)
                    }
                    value={ currentTeam.get('selectedMember') }
                    onChange={(val) => {
                      setTeamProperty('selectedMember', val.value)
                    }}
                    clearable={false}
                  />
                  <div
                    className={ "button" + (!!currentTeam.get('selectedMember') ? '' : ' disabled') }
                    onClick={() => {
                      if (!!currentTeam.get('selectedMember')) {
                        startEditingTeam()
                        addMember(currentTeam.get('selectedMember'))
                      }
                    }}
                  >
                    Add member
                  </div>
                </div>
              </td>
              <td className="add-object-label" colSpan="3" width="50%">
                Search by username, full name or email address
              </td>
            </tbody>
          </table>
      </div>
    ) 

    const selectedTeamContainer = !teams.size ? null : (
      <div id="selected-team" class="selected-tab">
        { selectedTeam }
      </div>
    )

    return (
      <div id="teams-section" className="section">
        <div className="section-header">
          <h1>Add team members to your expedition</h1>
          <p>process breadcrumbs</p>
        </div>
        <p className="intro">
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien. Praesent bibendum leo vitae felis pellentesque, sit amet mattis nisi mattis.
        </p>

        <ul id="teams-tabs" class="tabs">
          { teamTabs }
          <li className="team-name add" onClick={() => { addTeam() }}>+</li>
        </ul>
        { selectedTeamContainer }

        <Link to={'/admin/' + currentProjectID + '/' + expedition.get('id') }>
          <div className="button hero">
            Finalize your expedition!
          </div>
        </Link>
      </div>
    )    
  }
}

NewTeamsSection.propTypes = {
}

export default NewTeamsSection
