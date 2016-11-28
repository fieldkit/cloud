import React, {PropTypes, Component} from 'react'
import ViewportMercator from 'viewport-mercator-project'

// import THREE from '../react-three-renderer/node_modules/three'
// import React3 from '../react-three-renderer'

import autobind from 'autobind-decorator'

export default class WebGLOverlay extends Component {
  constructor (props) {
    super(props)

    const paths = {
      ambitGeo: []
    }

    const particles = {
      members: [],
      sightings: {
        count: 1000,
        // position: new THREE.BufferAttribute(new Float32Array(1000 * 3), 3),
        // color: new THREE.BufferAttribute(new Float32Array(1000 * 4), 4),
        // index: new THREE.BufferAttribute(new Uint16Array(1000 * 1), 1),
        data: [],
        vertexShader: [
          'attribute vec4 color;',
          'varying vec4 vColor;',
          'void main() {',
          '    vColor = color;',
          '    vec4 mvPosition = modelViewMatrix * vec4( position.xy, 0.0 , 1.0 );',
          '    gl_PointSize = float( position.z );',
          '    gl_Position = projectionMatrix * mvPosition;',
          '}'
        ].join('\n'),
        fragmentShader: [
          'varying vec4 vColor;',
          'uniform sampler2D texture;',
          'vec4 vFragColor;',
          'void main() {',
          '    vFragColor = vColor * texture2D( texture, gl_PointCoord );',
          '    if (vFragColor.w > 0.25) {',
          '      gl_FragColor = vFragColor;',
          '    } else {',
          '      discard;',
          '    }',
          '}'
        ].join('\n')
      },
      pictures360: {
        count: 1000,
        // position: new THREE.BufferAttribute(new Float32Array(1000 * 3), 3),
        // color: new THREE.BufferAttribute(new Float32Array(1000 * 4), 4),
        // index: new THREE.BufferAttribute(new Uint16Array(1000 * 1), 1),
        data: [],
        vertexShader: [
          'attribute vec4 color;',
          'varying vec4 vColor;',
          'void main() {',
          '    vColor = color;',
          '    vec4 mvPosition = modelViewMatrix * vec4( position.xy, 0.0 , 1.0 );',
          '    gl_PointSize = 20.0;',
          '    gl_Position = projectionMatrix * mvPosition;',
          '}'
        ].join('\n'),
        fragmentShader: [
          'varying vec4 vColor;',
          'uniform sampler2D texture;',
          'vec4 vFragColor;',
          'void main() {',
          '    vFragColor = vColor * texture2D( texture, gl_PointCoord );',
          '    if (vFragColor.w > 0.35) {',
          '      gl_FragColor = vFragColor;',
          '    } else {',
          '      discard;',
          '    }',
          '}'
        ].join('\n')
      }
    }

    for (var k in particles) {
      for (var i = 0; i < particles[k].count; i++) {
        particles[k].index.array[i] = i
      }
    }

    this.state = {
      paths,
      particles,
      render () {},
      // sightingTexture: new THREE.TextureLoader().load('/static/img/sighting.png'),
      // picture360Texture: new THREE.TextureLoader().load('/static/img/picture360.png'),
      mousePosition: [0, 0],
      hoveredPicture: -1
    }
  }

  componentWillReceiveProps (nextProps) {
    const { unproject } = ViewportMercator(nextProps)
    const render = nextProps.redraw({ unproject })

    if (!render) {
      this.setState({
        ...this.state,
        particles: null
      })
      return
    }

    const { particles, paths } = render(this.state.particles, this.state.paths)

    var hoveredPicture = -1
    if (this.state.mousePosition[0] > 0 || this.state.mousePosition[1] > 0) {
      const positions = particles.pictures360.position.array
      for (var i = 0; i < particles.pictures360.count; i ++) {
        const pos = [positions[i * 3 + 0], positions[i * 3 + 1]]
        const dist = Math.sqrt(Math.pow(this.state.mousePosition[0] - pos[0], 2) + Math.pow(this.state.mousePosition[1] - pos[1], 2))
        if (dist < 15) {
          particles.pictures360.color.array[i * 4 + 0] = 1
          particles.pictures360.color.array[i * 4 + 1] = 1
          particles.pictures360.color.array[i * 4 + 2] = 1
          particles.pictures360.color.array[i * 4 + 3] = 1
          hoveredPicture = i
          particles.pictures360.color.needsUpdate = true
          break
        }
        if (particles.pictures360.color.array[i * 4 + 3] === 0) break
      }
    }

    this.setState({
      ...this.state,
      particles,
      paths,
      render,
      hoveredPicture
    })
  }

  @autobind
  onMouseMove (event) {
    this.setState({
      ...this.state,
      mousePosition: [event.pageX, event.pageY]
    })
  }

  @autobind
  onMouseOut (event) {
    this.setState({
      ...this.state,
      mousePosition: [0, 0]
    })
  }

  @autobind
  onClick (event) {
    const { show360Picture } = this.props
    // const { particles } = this.state
    if (this.state.hoveredPicture > -1) {
      show360Picture(this.state.particles.pictures360.data[this.state.hoveredPicture])
    }
  }

  render () {
    const { project } = ViewportMercator(this.props)
    const { width, height, longitude, latitude } = this.props
    const { particles, paths } = this.state

    const point = project([longitude, latitude])
    const startPoint = project([this.state.longitude, this.state.latitude])
    const left = point[0] - startPoint[0]
    const top = 0 - (point[1] - startPoint[1])
    const cameraProps = {
      left: 0,
      right: width,
      top: 0,
      bottom: height,
      near: 1,
      far: 5000,
      // position: new THREE.Vector3(left, top, 600),
      // lookAt: new THREE.Vector3(left, top, 0)
    }

    const memberMarkers = particles.members
      .map((m, i) => {
        var x = Math.round((m.position[0] - 27 / 2) * 10) / 10
        var y = Math.round((m.position[1] - 34) * 10) / 10
        return (
          <div
            key={m.name}
            className={'member-marker'}
            style={{
              left: x,
              top: y
            }}
          >
            <img
              src="/static/img/member.svg"
              width={27}
              height={32}
            />
            <span>
              {m.name.charAt(0).toUpperCase()}
            </span>
          </div>
        )
      })

    const sightingLabels = particles.sightings.data
      .map((p, i) => {
        var x = particles.sightings.position.array[i * 3 + 0]
        var y = particles.sightings.position.array[i * 3 + 1]
        if (x >= window.innerWidth / 3 && x < 2 * window.innerWidth / 3 && y >= window.innerHeight / 3 && y < 2 * window.innerHeight / 3) {
        // if ((currentDate.getTime() - p.date.getTime()) < 200000 && (currentDate.getTime() - p.date.getTime()) > -200000) {
          return (
            <div
              key={i}
              className={'sighting-label'}
              style={{
                left: x,
                top: y
              }}
            >
              <div
                className="arrow-box"
              >
                {p.count + ' ' + p.name}
              </div>
            </div>
          )
        } else {
          return null
        }
      })

    const lines = paths.ambitGeo.map((p, i) => {

      return (
        <line key={i}>
          <geometry
            vertices={p.vertices}
            dynamic={true}
          >
          </geometry>
          <lineBasicMaterial
            linewidth={2}
            opacity={0.7}
            transparent={true}
            color={p.color}
          >
          </lineBasicMaterial>
        </line>
      )
    })

    return (
      <div>
        <div
          className={'hitbox' + (this.state.hoveredPicture > -1 ? ' hover' : '')}
          onMouseMove={this.onMouseMove}
          onMouseOut={this.onMouseOut}
          onClick={this.onClick}
        >
        </div>
        <div
          id="html-renderer"
          style={{
            position: 'absolute',
            width: '100%',
            height: '100%'
          }}
        >
         {sightingLabels}
         {memberMarkers}
        </div>
        <div id="three-renderer"
        >
          {/* 
          <React3
            mainCamera="camera"
            width={width}
            height={height}
            onAnimate={this._onAnimate}
            alpha={true}
            antialias={true}
          >
            <scene>
              <orthographicCamera
                name="camera"
                {... cameraProps}
              />
              { lines }
              { particles &&
                <points>
                  <bufferGeometry
                    position={particles.sightings.position}
                    index={particles.sightings.index}
                    color={particles.sightings.color}
                  />
                  <shaderMaterial
                    vertexShader={particles.sightings.vertexShader}
                    fragmentShader={particles.sightings.fragmentShader}
                    uniforms={
                      {texture: { type: 't', value: this.state.sightingTexture }}
                    }
                  >
                  </shaderMaterial>
                </points>
              }
              { particles &&
                <points>
                  <bufferGeometry
                    position={particles.pictures360.position}
                    index={particles.pictures360.index}
                    color={particles.pictures360.color}
                  />
                  <shaderMaterial
                    vertexShader={particles.pictures360.vertexShader}
                    fragmentShader={particles.pictures360.fragmentShader}
                    uniforms={
                      {texture: { type: 't', value: this.state.picture360Texture }}
                    }
                  >
                  </shaderMaterial>
                </points>
              }
            </scene>
          </React3>
          */ }
        </div>
      </div>
    )
  }

}

WebGLOverlay.propTypes = {
  show360Picture: PropTypes.func.isRequired,
  currentDate: PropTypes.object,
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
  latitude: PropTypes.number.isRequired,
  longitude: PropTypes.number.isRequired,
  zoom: PropTypes.number.isRequired,
  redraw: PropTypes.func.isRequired,
  isDragging: PropTypes.bool.isRequired
}
