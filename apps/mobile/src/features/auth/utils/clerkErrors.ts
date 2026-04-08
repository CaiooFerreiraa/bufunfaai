function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

function getFieldMessage(source: Record<string, unknown>, fieldName: string): string {
  const fields = source.fields;
  if (!isRecord(fields)) {
    return '';
  }

  const fieldValue = fields[fieldName];
  if (!isRecord(fieldValue)) {
    return '';
  }

  return typeof fieldValue.message === 'string' ? fieldValue.message : '';
}

function getGlobalMessage(source: Record<string, unknown>): string {
  const global = source.global;
  if (!Array.isArray(global) || global.length === 0) {
    return '';
  }

  const firstError = global[0];
  if (!isRecord(firstError)) {
    return '';
  }

  return typeof firstError.message === 'string' ? firstError.message : '';
}

export function getClerkErrorMessage(
  errors: unknown,
  fieldNames: readonly string[],
  fallback: string,
): string {
  if (!isRecord(errors)) {
    return fallback;
  }

  for (const fieldName of fieldNames) {
    const message = getFieldMessage(errors, fieldName);
    if (message) {
      return message;
    }
  }

  const globalMessage = getGlobalMessage(errors);
  if (globalMessage) {
    return globalMessage;
  }

  if (typeof errors.message === 'string' && errors.message.trim() !== '') {
    return errors.message;
  }

  return fallback;
}
