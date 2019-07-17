/**
 * Regular expression to identify valid username.
 */
export const VALID_USERNAME_REGEXP = /^[a-zA-Z0-9](?:[a-zA-Z0-9]|[-.](?=[a-zA-Z0-9]))*$/.source

/** Maximum allowed length for a username. */
export const USERNAME_MAX_LENGTH = 255

/** Maximum allowed length for a user display name. */
export const USER_DISPLAY_NAME_MAX_LENGTH = 255

/** Returns the URL path to a user's public profile */
export function userURL(username: string): string {
    return `/users/${username}`
}
