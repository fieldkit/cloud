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

export function errorsFor(errors: ?APIErrors, key: string): ?string {
  if (errors && errors.meta[`response.${key}`]) {
    return errors.meta[`response.${key}`];
  }
}

export function errorClass(errors: ?APIErrors, key: string): string {
  return errorsFor(errors, key) ? 'with-errors' : '';
}

export function slugify(input: string): string {
  return slug(input).toLowerCase();
}

export function joinPath(basePath: string, ...parts: string[]): string {
  const nextParts = parts.join('/');
  if (basePath.endsWith('/')) {
    basePath = basePath.substr(0, basePath.length - 1);
  }
  return `${basePath}/${nextParts}`;
}
