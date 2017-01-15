import React, {PropTypes, Component} from 'react'
import ViewportMercator from 'viewport-mercator-project'
import THREE from '../../vendors/react-three-renderer/node_modules/three'
import React3 from '../../vendors/react-three-renderer'
import sightingTexture from '../../img/sighting.png'

export default class WebGLOverlay extends Component {
  constructor (props) {
    super(props)

    const paths = {
      ambitGeo: []
    }

    const particles = {
      sightings: {
        count: 1000,
        position: new THREE.BufferAttribute(new Float32Array(1000 * 3), 3),
        color: new THREE.BufferAttribute(new Float32Array(1000 * 4), 4),
        index: new THREE.BufferAttribute(new Uint16Array(1000 * 1), 1),
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
      }
    }

    for (var k in particles) {
      for (var i = 0; i < particles[k].count; i++) {
        particles[k].index.array[i] = i
      }
    }

    this.state = {
      initialRender: false,
      paths,
      particles,
      render () {},
      sightingTexture: new THREE.TextureLoader().load(sightingTexture),
      mousePosition: [0, 0],
    }
  }

  componentWillReceiveProps (nextProps) {

    if (
      this.props.latitude !== nextProps.latitude ||
      this.props.longitude !== nextProps.longitude ||
      this.props.zoom !== nextProps.zoom ||
      this.props.width !== nextProps.width ||
      this.props.height !== nextProps.height
      ) {
      const { unproject } = ViewportMercator(nextProps)
      const render = nextProps.redraw({ unproject })
      const { particles, paths } = render(this.state.particles, this.state.paths)
      this.setState({
        ...this.state,
        particles,
        paths,
        render
      })
    }
  }

  // shouldComponentUpdate (nextProps) {
  //   return !this.state.initialRender
  // }

  // componentWillUpdate (nextProps) {
  // }  

  // componentDidUpdate () {
  //   this.setState({
  //     ...this.state,
  //     initialRender: true
  //   })
  // }

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
      position: new THREE.Vector3(left, top, 600),
      lookAt: new THREE.Vector3(left, top, 0)
    }

    // const sightingLabels = particles.sightings.data
    //   .map((p, i) => {
    //     var x = particles.sightings.position.array[i * 3 + 0]
    //     var y = particles.sightings.position.array[i * 3 + 1]
    //     if (x >= window.innerWidth / 3 && x < 2 * window.innerWidth / 3 && y >= window.innerHeight / 3 && y < 2 * window.innerHeight / 3) {
    //       return (
    //         <div
    //           key={i}
    //           className={'sighting-label'}
    //           style={{
    //             left: x,
    //             top: y
    //           }}
    //         >
    //           <div
    //             className="arrow-box"
    //           >
    //             {p.count + ' ' + p.name}
    //           </div>
    //         </div>
    //       )
    //     } else {
    //       return null
    //     }
    //   })

    // console.log('wooop', paths.ambitGeo)

    return (
      <div>
        <div
          className={'hitbox'}
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
         {/*sightingLabels*/}
         {/*memberMarkers*/}
        </div>
        <div id="three-renderer"
        >
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
                { ...cameraProps }
              />
              {
                paths &&
                // paths.ambitGeo.map((p, i) => {
                  // return (
                    <line>
                      <geometry
                        vertices={paths.ambitGeo}
                        dynamic={true}
                      >
                      </geometry>
                      <lineBasicMaterial
                        linewidth={10}
                        opacity={0.7}
                        transparent={false}
                        color={new THREE.Color('#ffffff')}
                      >
                      </lineBasicMaterial>
                    </line>
                  // )
                // }) 
              }
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
            </scene>
          </React3>
        </div>
      </div>
    )
  }

}

WebGLOverlay.propTypes = {
}
