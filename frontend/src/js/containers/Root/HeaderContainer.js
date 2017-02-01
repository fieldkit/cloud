
import { connect } from 'react-redux'
import Header from '../../components/Root/Header'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  const currentExpeditionID = state.expeditions.get('currentExpedition')
  const expeditionName = state.expeditions.getIn(['expeditions', currentExpeditionID, 'name'])

  return {
    expeditionName
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    saveChangesAndResume () {
      return dispatch(actions.saveChangesAndResume())
    },
    cancelAction () {
      return dispatch(actions.cancelAction())
    }
  }
}

const HeaderContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Header)

export default HeaderContainer
