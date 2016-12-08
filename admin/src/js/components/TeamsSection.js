import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'Immutable'

class TeamsSection extends React.Component {
  // constructor (props) {
    // super(props)
    // this.state = {
    //   ...props
    // }
  // }


  // @autobind
  // resetCurrentTeam () {
  //   this.setState({
  //     ...this.state,
  //     currentTeam: {
  //       name: '',
  //       editing: false,
  //       description: '',
  //       members: []
  //     }
  //   })
  // }

  // @autobind
  // selectMember (i) {
  //   this.setState({
  //     ...this.state,
  //     selectedMemberIndex: i
  //   })
  // }

  // @autobind
  // unselectMember (i) {
  //   this.setState({
  //     ...this.state,
  //     selectedMemberIndex: -1
  //   })
  // }

  // @autobind
  // onContentEdited (e) {
  //   const { updateExpedition } = this.props
  //   updateExpedition(this.state.expedition)
  // }

  // @autobind
  // onContentChange (e) {
  //   this.setState({
  //     ...this.state,
  //     expedition: this.state.expedition.setIn(['teams', this.state.selectedTeamIndex, 'description'], e.target.value)
  //   })
  // }

  // @autobind
  // resetChanges () {
  //   const { updateExpedition } = this.props
  //   this.setState({
  //     ...this.state,
  //     expedition: this.state.initialExpedition
  //   })
  //   updateExpedition(this.state.initialExpedition)
  // }

  // componentWillReceiveProps (nextProps) {
  //   this.setState(nextProps)
  // }

  render () {

    const { 
      expedition,
      teams,
      members,
      currentTeam,
      currentMember,
      editedTeam,
      setCurrentTeam,
      setCurrentMember,
      addTeam,
      startEditingTeam,
      stopEditingTeam,
      setTeamProperty,
      saveChangesToTeam,
      clearChangesToTeam
    } = this.props

    const teamTabs = teams
      .map((t, i) => {
        let className = 'team-name '
        if (t === currentTeam) className += 'editable active '
        if (currentTeam.get('name') === '') className += 'required '

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
            return <li key={j}>{d}</li>
          })

        if (!!currentMember && i !== currentMember.get('id')) {
          return (
            <tr key={i} onClick={() => {
              this.selectMember(i)
            }}>
              <td className="name">{m.get('name')}</td>
              <td className="role">{m.get('role')}</td>
              <td className="inputs">
                <ul>
                  { inputs }
                </ul>
              </td>
              <td className="activity">
                <svg></svg>
              </td>
              <td className="edit">
                <img src="/src/img/icon-edit-small"/>
              </td>
              <td className="remove">  
                <img src="/src/img/icon-remove-small"/>
              </td>
            </tr>
          )
        } else {
          return (
            <tr key={i} onClick={() => {
              this.unselectMember()
            }}>
              <td className="name" colSpan="6" width="100%">
                <h2>{m.get('name')}</h2>
              </td>
            </tr>
          )
        }
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
            !!currentTeam.find((val, key) => {
              return key !== 'status' && editedTeam.get(key) !== val
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
                  // if (!!currentTeam.get('description')) return currentTeam.get('description')
                  // else if (currentTeam.get('status') === 'new') return 'Enter a description for this new team.'
                  // else return ''
                })()}
                disabled={false}
                onClick={(e) => {
                  if (this.props.currentTeam.get('description') === 'Enter a description') setTeamProperty('description', '')
                  startEditingTeam()
                }}
                onBlur={(e) => {
                  // console.log('aga', currentTeam.get('description'), currentTeam.get('id'))
                  if (this.props.currentTeam.get('description') === '') setTeamProperty('description', 'Enter a description')
                  stopEditingTeam()
                }}
                onChange={(e) => { 
                  setTeamProperty('description', e.target.value)
                  // console.log(currentTeam.get('description'), currentTeam.get('id'))
                }}
              />
            </div>
            <svg className="activity column">
            </svg>
          </div>
          <table className="members-list">
            <tr>
              <th className="name">Name</th>
              <th className="role">Role</th>
              <th className="inputs">Inputs</th>
              <th className="activity">Activity</th>
              <th className="edit"></th>
              <th className="remove"></th>
            </tr>
            { teamMembers }
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
          {/*sectionActions*/}
          <h1>Teams</h1>
        </div>
        <p className="intro">
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien. Praesent bibendum leo vitae felis pellentesque, sit amet mattis nisi mattis.
        </p>

        <ul id="teams-tabs" class="tabs">
          { teamTabs }
          <li className="team-name add" onClick={() => { addTeam() }}>+</li>
        </ul>
        { selectedTeamContainer }
        {/*sectionActions*/}
      </div>
    )

      ////////
      ////////
      ////////
      ////////
      ////////
      ////////
      ////////



    return null
    // const { selectedTeamIndex, selectedMemberIndex, expedition, initialExpedition } = this.state
    // const selectedTeamIndex = this.state.get('selectedTeamIndex')
    // const selectedMemberIndex = this.state.get('selectedMemberIndex')
    // const expedition = this.state.get('expedition')

    

   

    // const sectionActions = (
    //   <ul className="section-actions">
    //     <li>
    //       { expedition.get('updating') &&
    //         (
    //           <div class="status">
    //             <span className="spinning-wheel-container"><div className="spinning-wheel"></div></span>
    //             Saving Changes
    //           </div>
    //         )
    //       }
    //       { !expedition.get('updating') &&
    //         initialExpedition !== this.props.expedition &&
    //         (
    //           <div class="status">
    //             Changes saved
    //           </div>
    //         )
    //       }
    //     </li>
    //     <li>
    //       <div class={'button primary ' + (initialExpedition === expedition ? 'disabled' : '')} onClick={this.resetChanges}>
    //         Reset Changes
    //       </div>
    //     </li>
    //   </ul>
    // )

    
  }
}

TeamsSection.propTypes = {
  expedition: PropTypes.object.isRequired,
  updateExpedition: PropTypes.func.isRequired
}

export default TeamsSection