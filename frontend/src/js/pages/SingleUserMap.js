// @flow weak

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Redirect } from 'react-router';

import UserSession from '../api/session';

import { generatePointDecorator } from '../common/utilities';

import MapContainer from '../components/MapContainer';

import '../../css/map.css';

function getMapLocationFromQueryString() {
    const params = new URLSearchParams(window.location.search);
    const value = params.get('center') || "";
    const center = value.split(",").map(v => Number(v));
    if (center.length < 2) {
        return null;
    }
    return center;
}

function getDefaultMapLocation() {
    const fromQuery = getMapLocationFromQueryString();
    if (fromQuery != null) {
        return fromQuery;
    }

    return [-118.2688137, 34.0309388, 14];
}

class DownloadDataPanel extends React.Component {
    render() {
        const { onDownload, onLogout } = this.props;

        return (
            <div className="download-data-panel">
                <div className="download-data-body">
                    <h4>Instructions Header Placeholder</h4>

                    <ol>
                        <li>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</li>
                        <li>Nunc hendrerit scelerisque semper. Donec pharetra nibh eu dui convallis, eget sagittis nunc pellentesque.</li>
                    </ol>

                    <button onClick={ onDownload } className="download">Download Data</button>

                    <button onClick={ onLogout } className="logout">Logout</button>
                </div>
            </div>
        );
    }
}

class SingleUserMap extends Component {
    static contextTypes = {
        router: PropTypes.shape({
            history: PropTypes.shape({
                push: PropTypes.func.isRequired,
                replace: PropTypes.func.isRequired,
            }).isRequired
        }).isRequired
    };

    constructor() {
        super();

        this.state = {
            focus: { center: getDefaultMapLocation() }
        };
    }

    // Not sure how I feel about this. I'm open to suggestions.
    componentWillMount() {
        // $FlowFixMe
        document.body.style.overflow = "hidden";
    }

    componentWillUnmount() {
        // $FlowFixMe
        document.body.style.overflow = null;
    }

    onUserActivity(map) {
        this.updateCoordinates(map);
    }

    updateCoordinates(map) {
        const newZoom = map.getZoom();
        const mapCenter = map.getCenter();
        const newCenter = [ mapCenter.lng, mapCenter.lat, newZoom ];
        this.context.router.history.push({
            pathname: '/map',
            search: '?center=' + newCenter.join(","),
        });
    }

    loadMapFeatures(criteria) {
    }

    async onDownload() {
        console.log("Download");
    }

    async onLogout() {
        await new UserSession().logout();

        this.setState({});
    }

    render() {
        const session = new UserSession();
        if (!session.authenticated()) {
            return <Redirect to={ "/" } />;
        }

        const pointDecorator = generatePointDecorator('constant', 'constant');
        const visibleFeatures = {
            focus: this.state.focus,
            geojson: {
                features: []
            }
        };

        return (
            <div className="map page">
                <div className="header">
                    <div className="project">Project</div>
                </div>
                <div>
                    <MapContainer style={{ }} containerStyle={{ width: "100vw", height: "100vh" }} controls={false}
                        visibleFeatures={ visibleFeatures }
                        pointDecorator={ pointDecorator }
                        focusFeature={ () => { } }
                        focusSource={ () => { } }
                        onUserActivity={ this.onUserActivity.bind(this) }
                        loadMapFeatures={ this.loadMapFeatures.bind(this) }
                        onChangePlaybackMode={ () => { } }>
                        <DownloadDataPanel onDownload={ this.onDownload.bind(this )} onLogout={ this.onLogout.bind(this )} />
                    </MapContainer>
                </div>
            </div>
        );
    }
}

export default SingleUserMap;
