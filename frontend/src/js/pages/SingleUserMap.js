// @flow weak

import _ from 'lodash';
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Redirect } from 'react-router';

import { FkPromisedApi } from '../api/calls';
import UserSession from '../api/session';

import { GeoRectSet, GeoRect } from '../common/geo';
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
        return (
            <div className="download-data-panel">
                <div className="download-data-body">
                    <h4>Instructions Header Placeholder</h4>

                    <ol>
                        <li>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</li>
                        <li>Nunc hendrerit scelerisque semper. Donec pharetra nibh eu dui convallis, eget sagittis nunc pellentesque.</li>
                    </ol>

                    {this.renderDownloadButton()}
                </div>
            </div>
        );
    }

    renderDownloadButton() {
        const { summary, onDownload } = this.props;

        if (_.isObject(summary) && !_.isEmpty(summary.urls.csv)) {
            return (
                <a target="_blank" rel="noopener noreferrer" href={summary.urls.csv}>
                    <button onClick={ onDownload } className="download">Download Data</button>
                </a>
            );
        }

        return <div />;
    }
}

class MapFeatures {
    constructor() {
        this.loaded = [];
        this.sources = {};
        this.geometries = {};
        this.geometriesBySource = {};
    }

    query_(criteria) {
        const desired = new GeoRect(criteria);
        const loaded = new GeoRectSet(this.loaded);
        if (loaded.contains(desired)) {
            console.log("query: ignoring");
            return Promise.resolve(false);
        }

        const loading = desired.enlarge(2);
        this.loaded.push(loading);

        console.log("query: querying");

        return FkPromisedApi.getMySimpleFeatures({
            ne: loading.ne,
            sw: loading.sw,
        }).then(data => {
            const incomingGeometries = _(data.geometries).map(g => {
                return [ g.id, g ];
            }).fromPairs().value();

            this.geometries = { ...this.geometries, ...incomingGeometries };

            const sourceIds = _.union(_(data.spatial).map(s => s.sourceId).value(), _(data.temporal).map(s => s.sourceId).uniq().value());
            return Promise.all(sourceIds.map(id => {
                console.log("promise source info", id);
                if (this.sources[id]) {
                    return this.sources[id];
                }
                return this.sources[id] = FkPromisedApi.getSource(id).then(source => {
                    return FkPromisedApi.getSourceSummary(id).then(summary => {
                        return { source, summary };
                    });
                });
            }));
        }).then(() => {
            return this.geometries;
        });
    }

    getSources() {
        console.log("getSources");
        return Promise.all(
            _(this.sources)
                .map((value, key) => {
                    console.log("source", key);
                    return value.then(ss => {
                        const geometries = _(this.geometries)
                              .map((value, key) => {
                                  return value;
                              })
                              .filter(g => g.sourceId === ss.source.id)
                              .value();
                        return {
                            source: ss.source,
                            summary: ss.summary,
                            geometries: geometries,
                        };
                    });
                })
                .value()).then((data) => {
                    console.log("sources data", data);
                    return data;
                }
        );
    }

    load(criteria) {
        console.log("load");
        return this.query_(criteria);
    }
};

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

        this.features = new MapFeatures();

        this.loadMapFeaturesThrottled = _.throttle(this.loadMapFeatures.bind(this), 500, { leading: true });

        this.state = {
            focus: { center: getDefaultMapLocation() }
        };
    }

    refreshSummary() {
        FkPromisedApi.getMySimpleSummary().then(summary => {
            this.setState({
                summary,
                focus: {
                    center: summary.center,
                },
            });
        });
    }

    loadMapFeatures(criteria) {
        return this.features.load(criteria).then(geometries => {
            if (geometries) {
                return this.features.getSources().then(sources => {
                    this.setState({
                        geometries: geometries,
                        sources: sources,
                    });
                });
            }
        }, () => {
            this.setState({});
        });
    }

    // Not sure how I feel about this. I'm open to suggestions.
    componentWillMount() {
        this.refreshSummary();

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

    async onDownload() {
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
            sources: this.state.sources,
            focus: this.state.focus
        };

        return (
            <div className="map page">
                <div className="header">
                    <div className="project">Project</div>
                    <button onClick={ this.onLogout.bind(this) } className="logout">Logout</button>
                </div>
                <div>
                    <MapContainer style={{ }} containerStyle={{ width: "100vw", height: "100vh" }} controls={false}
                        visibleFeatures={ visibleFeatures }
                        pointDecorator={ pointDecorator }
                        focusFeature={ () => { } }
                        focusSource={ () => { } }
                        onUserActivity={ this.onUserActivity.bind(this) }
                        loadMapFeatures={ this.loadMapFeaturesThrottled }
                        onChangePlaybackMode={ () => { } }>
                        <DownloadDataPanel onDownload={ this.onDownload.bind(this )} summary={this.state.summary} />
                    </MapContainer>
                </div>
            </div>
        );
    }
}

export default SingleUserMap;
