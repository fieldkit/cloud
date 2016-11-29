import React, {PropTypes} from 'react'
import { Link } from 'react-router'

class TeamsSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      selectedTeamIndex: 0
    }
  }

  render () {

    const { expedition } = this.props
    const { selectedTeamIndex } = this.state

    const tabLabels = expedition.get('teams')
      .map((t, i) => {
        return (
          <li className={'team-name ' + (i === selectedTeamIndex ? 'active' : '')} key={i}>{t.get('name')}</li>
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

        return (
          <tr key={i}>
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
            <td className="data">
              <img src="/src/img/icon-data-small"/>
            </td>
            <td className="edit">
              <img src="/src/img/icon-edit-small"/>
            </td>
            <td className="remove">  
              <img src="/src/img/icon-remove-small"/>
            </td>
          </tr>
        )
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
                <div className="column">
                  <h5>Description</h5>
                  <p className="description">
                    {t.get('description')}
                  </p>
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
                  <th className="data"></th>
                  <th className="edit"></th>
                  <th className="remove"></th>
                </tr>
                {members}
              </table>
          </div>
        )
      })

    return (
      <div id="teams-section" className="section">
        <h1>Teams</h1>
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

      </div>
    )
  }
}

TeamsSection.propTypes = {
  expedition: PropTypes.object.isRequired
}

export default TeamsSection