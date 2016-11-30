import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'Immutable'

class TeamsSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      selectedTeamIndex: 0,
      selectedMemberIndex: -1,
      initialExpedition: props.expedition,
      expedition: props.expedition
    }
  }

  @autobind
  selectTeam (i) {
    this.setState(
      this.state.set('selectedTeamIndex', i)
    )
  }

  @autobind
  selectMember (i) {
    this.setState(
      this.state.set('selectedMemberIndex', i)
    )
  }

  @autobind
  unselectMember (i) {
    this.setState(
      this.state.set('selectedMemberIndex', -1)
    )
  }

  @autobind
  onContentEdited (e) {
    const { updateExpedition } = this.props
    updateExpedition(this.state.expedition)
  }

  @autobind
  onContentChange (e) {
    this.setState({
      ...this.state,
      expedition: this.state.expedition.setIn(['teams', this.state.selectedTeamIndex, 'description'], e.target.value)
    })
  }

  componentWillReceiveProps (nextProps) {

    const { expedition } = nextProps
    this.setState({
      ...this.state,
      expedition
    })
  }

  render () {
    const { selectedTeamIndex, selectedMemberIndex, expedition, initialExpedition } = this.state
    // const selectedTeamIndex = this.state.get('selectedTeamIndex')
    // const selectedMemberIndex = this.state.get('selectedMemberIndex')
    // const expedition = this.state.get('expedition')

    console.log('woop', this.state.expedition.get('loading'))

    const tabLabels = expedition.get('teams')
      .map((t, i) => {
        return (
          <li 
            className={'team-name ' + (i === selectedTeamIndex ? 'active' : '')}
            key={i}
            onClick={() => {
              this.selectTeam(i)
            }}
          >
            {t.get('name')}
          </li>
        )
      })

    const members = expedition.get('teams')
      .find((t, i) => {
        return i === selectedTeamIndex
      })
      .get('members')
      .map((m, i) => {
        const inputs = m.get('inputs')
          .map((d, j) => {
            return <li key={j}>{d}</li>
          })

        if (i !== selectedMemberIndex) {
          return (
            <tr key={i} onClick={() => {
              this.selectMember(i)
            }}>
              <td className="name">{m.get('name')}</td>
              <td className="role">{m.get('role')}</td>
              <td className="inputs">
                <ul>
                  {inputs}
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

    const selectedTeam = expedition.get('teams')
      .filter((t, i) => {
        return i === selectedTeamIndex
      })
      .map((t, i) => {
        return (
          <div className="team" key={i}>
              <div className="actions">
                <div className="button secondary">
                  Remove team
                </div>
              </div>
              <div className="header">
                <div className="column description">
                  <h5>Description</h5>
                  <ContentEditable
                    html={t.get('description')}
                    disabled={false}
                    onBlur={this.onContentEdited}
                    onChange={this.onContentChange}
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
                {members}
              </table>
          </div>
        )
      })

    const actions = () => {
      return (
        <ul className="actions">
          <li>
            { expedition.get('updating') &&
              (
                <div class="status">
                  <span className="spinning-wheel"></span> Saving Changes
                </div>
              )
            }
            { !expedition.get('updating') &&
              initialExpedition !== expedition &&
              (
                <div class="status">
                  Changes saved
                </div>
              )
            }
          </li>
          <li>
            <div class={'button primary ' + (initialExpedition === expedition ? 'diabled' : '')}>
              Reset Changes
            </div>
          </li>
        </ul>
      )
    }

    return (
      <div id="teams-section" className="section">
        <div className="section-header">
          {actions()}
          <h1>Teams</h1>
        </div>
        <p className="intro">
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien. Praesent bibendum leo vitae felis pellentesque, sit amet mattis nisi mattis.
        </p>

        <ul id="teams-tabs" class="tabs">
          {tabLabels}
          <li className="team-name">+</li>
        </ul>

        <div id="selected-team" class="selected-tab">
          {selectedTeam}
        </div>
        {actions()}
      </div>
    )
  }
}

TeamsSection.propTypes = {
  expedition: PropTypes.object.isRequired,
  updateExpedition: PropTypes.func.isRequired
}

export default TeamsSection