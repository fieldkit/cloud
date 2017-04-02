import React, {PropTypes, Component} from 'react'
import React3 from 'react-three-renderer'
import ViewportMercator from 'viewport-mercator-project'
import { Vector3, BufferAttribute, Color, TextureLoader } from 'three'

import iconSensorReading from '../../../img/sighting.png'

export default class WebGLOverlay extends Component {
  constructor (props) {
    super(props)

    // here you can add other feature types and assign them a particle texture
    const bufferGeometries = {
      particles: {
        'Feature': {
          count: 1000,
          position: new BufferAttribute(new Float32Array(1000 * 3), 3),
          color: new BufferAttribute(new Float32Array(1000 * 4), 4),
          index: new BufferAttribute(new Uint16Array(1000 * 1), 1),
          texture: new TextureLoader().load(('/' + iconSensorReading))
        }
      }
    }
    for (let k in bufferGeometries) {
      for (let l in bufferGeometries[k]) {
        for (let i = 0; i < bufferGeometries[k][l].count; i++) {
          bufferGeometries[k][l].index.array[i] = i
        }
      }
    }

    const shaders = {
      particles: {
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
    
    this.state = {
      shaders,
      bufferGeometries
    }
  }

  componentWillReceiveProps (nextProps) {
    const { bufferGeometries } = this.state
    const {
      particles
    } = nextProps
    

    Object.keys(particles).forEach(type => {
      bufferGeometries.particles[type].position.array = particles[type].position
      bufferGeometries.particles[type].position.needsUpdate = true
      bufferGeometries.particles[type].color.array = particles[type].color
      bufferGeometries.particles[type].color.needsUpdate = true
    })

    this.setState({
      ...this.state,
      bufferGeometries
    })
  }

  render () {
    const { project } = ViewportMercator(this.props)
    const { 
      width,
      height,
      longitude,
      latitude,
      readingPath
    } = this.props

    const { 
      bufferGeometries,
      shaders
     } = this.state

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
      position: new Vector3(left, top, 600),
      lookAt: new Vector3(left, top, 0)
    }

    return (
      <div>
        <div id="three-renderer">
          <React3
            mainCamera="camera"
            width={ width }
            height={ height }
            alpha={ true }
            antialias={ true }
          >
            <scene>
              <orthographicCamera
                name="camera"
                { ...cameraProps }
              />
              { /*
                !!readingPath &&
                <line>
                  <geometry
                    vertices={readingPath}
                    dynamic={true}
                  >
                  </geometry>
                  <lineBasicMaterial
                    linewidth={10}
                    opacity={0.7}
                    transparent={false}
                    color={new Color('#ffffff')}
                  >
                  </lineBasicMaterial>
                </line>
                */
              }
              {
                Object.keys(bufferGeometries.particles).map(type => {
                  const particles = bufferGeometries.particles[type]
                  return (
                    <points key={ 'particles-' + type }>
                      <bufferGeometry
                        position={ particles.position }
                        color={ particles.color }
                        index={ particles.index }
                      />
                      <shaderMaterial
                        vertexShader={ shaders.particles.vertexShader}
                        fragmentShader={ shaders.particles.fragmentShader}
                        uniforms={{
                          texture: { 
                            type: 't',
                            value: particles.texture 
                          }
                        }}
                      >
                      </shaderMaterial>
                    </points>
                  )
                })
              }
            </scene>
          </React3>
        </div>
      </div>
    )
  }
}

WebGLOverlay.propTypes = {}
