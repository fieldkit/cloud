// @flow

import React, { Component } from 'react';

import { errorsFor } from '../../common/util';
import type { APIErrors } from '../../api/types';

type Props = {
  labelText: string;
  name: string;
  className: string;
  type?: string;
  ref?: string;
  value: string;
  errors: ?APIErrors;
  onChange?: (e: any) => void;
}

export class FormItem extends Component<void, $Exact<Props>, void> {
  render() {
    const {
      name,
      className,
      labelText,
      type,
      ref,
      value,
      errors,
      onChange
    } = this.props;

    return (
      <div className="form-group">
        <label htmlFor={name}>{labelText}</label>
        <input ref={ref} name={name} type={type || 'text'} value={value} onChange={onChange}  className={className}/>
        { errorsFor(errors, name) }
      </div>
    );
  }
}
