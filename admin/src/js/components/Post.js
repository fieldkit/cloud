
import React, {PropTypes} from 'react'
import { dateToString } from '../utils'
// import THREE from '../react-three-renderer/node_modules/three'
// import React3 from '../react-three-renderer'
import autobind from 'autobind-decorator'

class Post extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      width: 640,
      height: 360,
      isUserInteracting: false,
      onMouseDownMouseX: 0,
      onMouseDownMouseY: 0,
      lon: 0,
      onMouseDownLon: 0,
      lat: 0,
      onMouseDownLat: 0,
      phi: 0,
      theta: 0,
      fov: 75,
      closeMouseOver: false
    }
  }

  componentWillReceiveProps (nextProps) {
    const { data } = this.props

    if (data.type === '360Image') {

      const width = document.querySelector('div.post' + data.key).clientWidth * 0.85
      const height = window.innerHeight - 250

      this.setState({
        ...this.state,
        width,
        height
      })
    }
  }

  @autobind
  onMouseDown (event) {
    // cosole.log('down)
    const {lon, lat} = this.state
    const { clientX, clientY } = event
    this.setState({
      ...this.state,
      isUserInteracting: true,
      onMouseDownMouseX: clientX,
      onMouseDownMouseY: clientY,
      onMouseDownLon: lon,
      onMouseDownLat: lat
    })
    event.stopPropagation()
  }

  @autobind
  onMouseMove (event) {
    const { isUserInteracting, onMouseDownMouseX, onMouseDownLon, onMouseDownMouseY, onMouseDownLat } = this.state
    if (isUserInteracting === true) {
      const lon = (onMouseDownMouseX - event.clientX) * 0.1 + onMouseDownLon
      const lat = Math.max(-85, Math.min(85, ((event.clientY - onMouseDownMouseY) * 0.1 + onMouseDownLat)))
      // const phi = THREE.Math.degToRad(90 - lat)
      // const theta = THREE.Math.degToRad(lon)
      // const cameraTarget = new THREE.Vector3(
      //   500 * Math.sin(phi) * Math.cos(theta),
      //   500 * Math.cos(phi),
      //   500 * Math.sin(phi) * Math.sin(theta)
      // )

      this.setState({
        ...this.state,
        lon,
        lat,
        cameraTarget
      })

    }
  }

  @autobind
  onMouseUp (event) {
    this.setState({
      ...this.state,
      isUserInteracting: false
    })
    event.stopPropagation()
  }

  @autobind
  onClick (event) {
    event.stopPropagation()
  }

  @autobind
  onMouseOut (event) {
    this.setState({
      ...this.state,
      isUserInteracting: false
    })
    event.stopPropagation()
  }

  @autobind
  onMouseWheel (event) {
    // const { fov } = this.state
    // this.setState({
    //   ...this.state,
    //   fov: fov + event.deltaY * 0.05
    // })
  }

  render () {
    const { format, data, closeLightBox, show360Picture } = this.props
    const { width, height, sphereScale, fov } = this.state

    if (format === 'full') {
      var metaTypes = ['date']
      if (data.location) metaTypes.push('location')
      if (data.link) metaTypes.push('link')
      var meta = metaTypes.map((m, i) => {
        var value = (function () {
          if (m === 'date') {
            var d = new Date(data.date)
            return (<p>{dateToString(d, true)}</p>)
          } else if (format === 'full') {
            return (<img width="16" height="16" key={i} src={'/static/img/icon-' + m + '.png'}/>)
          }
        })()
        return (
          <div className={m} key={i}>{value}</div>
        )
      })

      const content = data => {
        switch (data.type) {
          case 'tweet':
            var images = data.images ? data.images.map((url, i) => {
              return <img src={url} key={i}/>
            }) : ''
            return (
              <div>
                <p className="message">"{data.content}"</p>
                {images}
              </div>
            )
          case 'image':
            return (
              <div>
                <img src={data.images[0]}/>
                <p>{data.content}</p>
              </div>
            )
          case '360Image':
            return (
              <div
                onMouseDown={this.onMouseDown}
                onMouseMove={this.onMouseMove}
                onMouseUp={this.onMouseUp}
                onMouseOut={this.onMouseOut}
                onWheel={this.onMouseWheel}
                onClick={this.onClick}
              >
                <div className="interaction-helper">
                  <div>
                    <img src="/static/img/icon-360.svg"/>
                    Drag image to rotate view
                  </div>
                </div>
                {/* 
                <React3
                  mainCamera="lightbox-camera"
                  width={width}
                  height={height}
                >
                  <scene>
                    <perspectiveCamera
                      name="lightbox-camera"
                      fov={fov}
                      aspect={width / height}
                      near={1}
                      far={1100}
                      position={new THREE.Vector3(0, 0, 0)}
                      lookAt={this.state.cameraTarget}
                    />
                    <mesh
                      scale={sphereScale}
                    >
                      <sphereGeometry
                        radius={500}
                        widthSegments={60}
                        heightSegments={40}
                      />
                      <meshBasicMaterial>
                        <texture url={data.images[0]}/>
                      </meshBasicMaterial>
                    </mesh>
                  </scene>
                </React3>
                */ }
              </div>
            )
          case 'audio':
            var soundCloudURL = 'https://w.soundcloud.com/player/?url='
            soundCloudURL += data.link
            soundCloudURL += '&color=F8D143&auto_play=false&hide_related=false&show_comments=true&show_user=true&show_reposts=false'

            return (
              <div>
                <h2>{data.title}</h2>
                <p>sound</p>
                <iframe src={soundCloudURL}></iframe>
              </div>
            )
          case 'blog':
            return (
              <div>
                <h2>{data.title}</h2>
                <p>{data.content}</p>
                <p><a href={data.link}>Read more...</a></p>
              </div>
            )
          default:
            return ''
        }
      }

      return (
        <div className={'post post' + data.key}>
          <div className="type">
            <img width="16" height="16" src={'/static/img/icon-' + data.type + '.png'}/>
          </div>
          <div className="content">
            {content(data)}
            {
              data.type !== '360Image' &&
              <div className="meta">
                {meta}
                <div className="share"><img width="16" height="16" src="/static/img/icon-share.png"/></div>
                <div className="separator"></div>
              </div>
            }
          </div>
          {
            data.type === '360Image' &&
            <div className="actions">
              <div
                className="button close"
                onMouseOver={() => {
                  this.setState({
                    ...this.state,
                    closeMouseOver: true
                  }) }
                }
                onMouseOut={() => {
                  this.setState({
                    ...this.state,
                    closeMouseOver: false
                  }) }
                }
                onClick={closeLightBox}
              >
                <img
                  width="16"
                  height="16"
                  src={'/static/img/icon-close' + (this.state.closeMouseOver ? '-hover' : '') + '.png'}
                />
              </div>
              <div
                className="button forward"
                onMouseOver={() => {
                  this.setState({
                    ...this.state,
                    forwardMouseOver: true
                  }) }
                }
                onMouseOut={() => {
                  this.setState({
                    ...this.state,
                    forwardMouseOver: false
                  }) }
                }
                onMouseDown={(event) => {
                  if (data.next) show360Picture(data.next)
                  event.stopPropagation()
                }}
              >
                <img
                  width="16"
                  height="16"
                  src={'/static/img/icon-forward' + (this.state.forwardMouseOver ? '-hover' : '') + '.png'}
                />
              </div>
              <div
                className="button backward"
                onMouseOver={() => {
                  this.setState({
                    ...this.state,
                    backwardMouseOver: true
                  }) }
                }
                onMouseOut={() => {
                  this.setState({
                    ...this.state,
                    backwardMouseOver: false
                  }) }
                }
                onMouseDown={(event) => {
                  if (data.previous) show360Picture(data.previous)
                  event.stopPropagation()
                }}
              >
                <img
                  width="16"
                  height="16"
                  src={'/static/img/icon-backward' + (this.state.backwardMouseOver ? '-hover' : '') + '.png'}
                />
              </div>
            </div>
          }
        </div>
      )
    }

    return ''
  }
}

Post.propTypes = {
  format: PropTypes.string.isRequired,
  data: PropTypes.object.isRequired,
  closeLightBox: PropTypes.func,
  show360Picture: PropTypes.func
}

export default Post


