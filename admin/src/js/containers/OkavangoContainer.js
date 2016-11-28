
import { connect } from 'react-redux'
import Okavango from '../components/Okavango'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  return {
    expedition: state.expeditions[state.selectedExpedition],
    children: ownProps.children,
    animate: state.animate,
    isFetching: state.isFetching,
    mapStateNeedsUpdate: state.mapStateNeedsUpdate,
    expeditionID: state.selectedExpedition,
    contentActive: state.contentActive,
    initialPage: state.initialPage,
    lightBoxActive: state.lightBoxActive,
    lightBoxPost: state.lightBoxPost ? formatPost(state.lightBoxPost, state.selectedExpedition) : null
  }
}

const formatPost = (p, expeditionID) => {
  var key = p.id
  var type = p.properties.FeatureType
  var date = new Date(p.properties.DateTime)
  var location = p.geometry.coordinates
  var author = p.properties.Member
  var title, content, images, link, dimensions, next, previous

  if (type === 'tweet') {
    if (expeditionID !== 'okavango_14') {
      content = p.properties.Text
      images = p.properties.Images.map(i => { return i.Url.replace('http://', 'https://') })
      link = p.properties.Url.replace('http://', 'https://')
    } else {
      content = p.properties.Tweet.text
    }
  }

  if (type === 'image') {
    if (expeditionID !== 'okavango_14') {
      content = p.properties.Notes
      images = [p.properties.Url.replace('http://', 'https://')]
      link = p.properties.InstagramID
      dimensions = p.properties.Dimensions
    } else {
      content = p.properties.Notes
      images = [p.properties.Url.replace('http://', 'https://')]
      link = p.properties.InstagramID
      dimensions = p.properties.Size
    }
    if (p.properties.Make === 'RICOH') {
      type = '360Image'
      next = p.properties.next
      previous = p.properties.previous
    }
  }

  if (type === 'blog') {
    title = p.properties.Title
    content = p.properties.Summary
    link = p.properties.Url.replace('http://', 'https://')
  }

  if (type === 'audio') {
    title = p.properties.Title
    link = p.properties.SoundCloudURL.replace('http://', 'https://')
  }

  return {
    key,
    type,
    title,
    content,
    images,
    link,
    date,
    location,
    author,
    dimensions,
    next,
    previous
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    fetchDay: (currentDate) => {
      return dispatch(actions.fetchDay(currentDate))
    },
    updateMap: (currentDate, coordinates, viewGeoBounds, zoom, expeditionID) => {
      return dispatch(actions.updateMap(currentDate, coordinates, viewGeoBounds, zoom, expeditionID))
    },
    setControl: (target, mode) => {
      return dispatch(actions.setControl(target, mode))
    },
    jumpTo: (date, expeditionID) => {
      return dispatch(actions.jumpTo(date, expeditionID))
    },
    setPage: () => {
      return dispatch(actions.setPage())
    },
    enableContent: () => {
      return dispatch(actions.enableContent())
    },
    show360Picture: (post) => {
      return dispatch(actions.show360Picture(post))
    },
    closeLightBox: () => {
      return dispatch(actions.closeLightBox())
    }
  }
}

const OkavangoContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Okavango)

export default OkavangoContainer
