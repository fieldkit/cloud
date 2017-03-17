// @flow
import React from 'react';
import slug from 'slug';

import type { APIErrors } from '../api/types';

export class BaseError {
  name: string;
  message: string;
  stack: ?string;

  constructor(message: string = 'Error') {
    this.name = this.constructor.name;
    this.message = message;
    if (typeof Error.captureStackTrace === 'function') {
      Error.captureStackTrace(this, this.constructor);
    } else {
      this.stack = (new Error(message)).stack;
    }
  }
}

export function errorsFor(errors: ?APIErrors, key: string) {
  if (errors && errors.meta[key]) {
    return errors.meta[key];
  }
}

export function slugify(input: string): string {
  return slug(input).toLowerCase();
}
