// @flow weak

import React, { Component } from 'react'

import fkLogo from '../../../img/logos/fieldkit-logo-red.svg'
import '../../../css/unauth.css'

export class Unauth extends Component {
    render() {
        return (
            <div className="unauth">
                <div className="contents">
                    { this.props.children }
                </div>
            </div>
        )
    }
}
