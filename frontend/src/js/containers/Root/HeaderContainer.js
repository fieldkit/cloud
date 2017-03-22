
import { connect } from 'react-redux'
import Header from '../../components/Root/Header'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  const currentExpeditionID = state.expeditions.get('currentExpedition')
  const expeditionName = state.expeditions.getIn(['expeditions', currentExpeditionID, 'name'])
  const currentPage = state.expeditions.get('currentPage')

  return {
    expeditionName,
    currentExpeditionID,
    currentPage
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    saveChangesAndResume () {
      return dispatch(actions.saveChangesAndResume())
    },
    cancelAction () {
      return dispatch(actions.cancelAction())
    },
    openExpeditionPanel () {
      return dispatch(actions.openExpeditionPanel())
    }
  }
}

const HeaderContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Header)

export default HeaderContainer
