// @flow weak

import _ from 'lodash';
import moment from 'moment';

import { API_HOST } from '../../secrets';

import React, { Component } from 'react';

const panelStyle: React.CSSProperties = {
    backgroundColor: '#f9f9f9',
    color: "#000",
    borderTopLeftRadius: 2,
    borderTopRightRadius: 2,
    borderBottomLeftRadius: 2,
    borderBottomRightRadius: 2,
    padding: '10px',
    // boxShadow: '0px 1px 4px rgba(0, 0, 0, .3)',
    // border: '1px solid rgba(0, 0, 0, 0.4)'
};

const containerStyle: React.CSSProperties = {
    position: 'absolute',
    width: '300px',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'column',
};

const sourcesTitleStyle: React.CSSProperties = {
    fontSize: '18px',
    fontWeight: 'bold',
    paddingBottom: '5px',
};

const sourceContainerStyle: React.CSSProperties = {
    paddingBottom: '10px',
};

const sourceNameStyle: React.CSSProperties = {
    fontWeight: 'bold',
};

const sourceFirstLineStyle: React.CSSProperties = {
    fontSize: '14px',
};

const sourceSecondLineStyle: React.CSSProperties = {
};

export default class FiltersPanel extends Component {
    props: {
        style: React.CSSProperties
    }

    render() {
        const { style, visibleFeatures } = this.props;
        const position = { top: 160, left: 20, bottom: 'auto', right: 'auto' };

        if (_.values(visibleFeatures.sources).length === 0) {
            return (<div></div>);
        }

        return (
            <div style={{ ...panelStyle, ...containerStyle, ...style, ...position }}>
                {this.renderSources(visibleFeatures.sources)}
            </div>
        );
    }

    renderSources(sources) {
        const { onShowSource, onShowFeature } = this.props;

        return (
            <div>
            <div style={sourcesTitleStyle}>Sources</div>
            {_.values(sources).filter(r => r.source && r.summary).map(r => {
                const source = r.source;
                const lastFeature = r.lastFeature;
                // const summary = r.summary;
                if (lastFeature == null) {
                    return (<div key={source.id}>Loading...</div>);
                }
                const lastFeatureDate = moment(new Date(lastFeature.properties.timestamp)).format('MMM Do YYYY, h:mm:ss a');
                const numberOfFeatures = source.numberOfFeatures;
                const featuresUrl = API_HOST + "/inputs/" + source.id + "/geojson";

                return (
                    <div key={source.id} style={sourceContainerStyle}>
                        <div style={sourceFirstLineStyle}>
                            <a style={{ ...sourceNameStyle, ...{ color: 'black' } }} href="#" onClick={() => onShowSource(r)}>{source.name}</a>
                            &nbsp;
                            <a style={{ color: 'black' }} target="_blank" href={featuresUrl}>{numberOfFeatures} features</a>
                        </div>
                        <div style={sourceSecondLineStyle}>
                            <a style={{ color: 'black' }} href="#" onClick={() => onShowFeature(lastFeature) }>{lastFeatureDate}</a>
                        </div>
                    </div>
                );
            })}
            </div>
        );
    }
}
