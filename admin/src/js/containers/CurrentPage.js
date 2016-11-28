
import { connect } from 'react-redux'
import Content from '../components/Content'
import { changeCurrentPage } from '../actions'


const mapStateToProps = (state, ownProps) => {
  return {
    currentPage: state.currentPage
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    // onNavClick: (target) => {
    //   dispatch({'type':'NAV', 'target':target})
    // }
  }
}

const CurrentPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(Content)


export default CurrentPage