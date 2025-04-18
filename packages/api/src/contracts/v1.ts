/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */

export interface paths {
    "/v1/identity/disable-two-factor": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /** Disable two-factor authentication */
        delete: operations["disable-two-factor"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/enable-two-factor": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Enable two-factor authentication after setup */
        post: operations["enable-two-factor"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/forgot-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Initiate password reset process */
        post: operations["forgot-password"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/me": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get current user session */
        get: operations["get-session"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/refresh-session": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Extend current session validity */
        post: operations["refresh-session"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/regenerate-qr-code": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Regenerate QR code for existing two-factor authentication setup */
        post: operations["regenerate-qr-code"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/reset-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Complete password reset process */
        post: operations["reset-password"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/sessions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Fetch all active sessions */
        get: operations["get-all-sessions"];
        put?: never;
        post?: never;
        /** Terminate all active sessions */
        delete: operations["delete-all-sessions"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/sessions/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /** Terminate session by id */
        delete: operations["delete-session"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/setup-two-factor": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Setup two-factor authentication */
        post: operations["setup-two-factor"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/sign-in": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Authenticate and create a new session */
        post: operations["sign-in"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/sign-out": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /** Terminate current session */
        delete: operations["sign-out"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/sign-up": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Create a new user account */
        post: operations["sign-up"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/update-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Update user password */
        post: operations["update-password"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/verify-email": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Confirm user email address */
        post: operations["verify-email"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/verify-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Verify current password for sensitive operations */
        post: operations["verify-password"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/identity/verify-two-factor": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Verify two-factor authentication code during sign-in */
        post: operations["verify-two-factor"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/payments": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Create payment */
        post: operations["create-payment"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/users/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get user */
        get: operations["get-user"];
        /** Update user */
        put: operations["update-user"];
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/users/{id}/image": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Update user profile image */
        post: operations["update-user-image"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        CreatePaymentRequestBody: {
            /** @description All monetary amounts should be provided in the minor unit. For example:
             *       - 1000 to charge 10 USD (or any other two-decimal currency)
             *       - 10 to charge 10 JPY (or any other zero-decimal currency) */
            Amount: number | string;
            /**
             * Format: currency
             * @description The currency code in ISO 4217 format.
             */
            Currency: string;
            Description: string;
            MerchantID: string;
            Metadata: {
                [key: string]: unknown;
            };
            Method: string;
            Provider: string;
        };
        DisableTwoFactorResponseBody: {
            /** @description Whether two-factor authentication is enabled */
            enabled: boolean;
        };
        EnableTwoFactorRequestBody: {
            /** @description The verification code to confirm setup */
            code: string;
        };
        EnableTwoFactorResponseBody: {
            /** @description Whether two-factor authentication is enabled */
            enabled: boolean;
        };
        Entity: {
            /** @description The entity's domain */
            domain?: string;
            /** @description The entity's ID */
            id: string;
            /** @description The entity's logo URL */
            logo?: string;
            /** @description The entity's name */
            name: string;
            /** @description The parent entity's ID */
            parentId?: string;
            /** @description The entity's slug */
            slug: string;
            /** @description The entity's status */
            status: string;
            /** @description The entity's type */
            type: string;
        };
        EntityRole: {
            /** @description The role's access to resources in the entity */
            access: {
                [key: string]: string[] | null;
            };
            /** @description The name of the user's role in the entity */
            name: string;
        };
        Error: {
            code?: components["schemas"]["ErrorCode"];
            errors: components["schemas"]["ErrorDetail"][] | null;
            message: string;
        };
        /**
         * @description The standardized error code
         * @enum {string}
         */
        ErrorCode: "AccountLocked" | "BackupCodeValidation" | "ConnectionNotFound" | "DuplicateItems" | "EmailExists" | "EmailNotVerified" | "EntityNotFound" | "FailedToVerifyTurnstileToken" | "InsufficientPermissions" | "InvalidBody" | "InvalidConnectionCredentials" | "InvalidCountry" | "InvalidCredentials" | "InvalidCurrency" | "InvalidCursor" | "InvalidDate" | "InvalidDateTime" | "InvalidEmail" | "InvalidFinancialAmount" | "InvalidHostname" | "InvalidIPv4" | "InvalidIPv6" | "InvalidImageFormat" | "InvalidName" | "InvalidOrExpiredToken" | "InvalidRefreshToken" | "InvalidTime" | "InvalidTurnstileToken" | "InvalidTwoFactorCode" | "InvalidUUID" | "InvalidValue" | "MissingLowercase" | "MissingNumber" | "MissingSpecial" | "MissingUppercase" | "PaymentNotFound" | "Required" | "TooLarge" | "TooLong" | "TooShort" | "TooSmall" | "TwoFactorAlreadyEnabled" | "TwoFactorLocked" | "TwoFactorNotEnabled" | "TwoFactorPending" | "Unauthenticated" | "Unknown" | "Unused" | "UserNotFound";
        ErrorDetail: {
            code: components["schemas"]["ErrorCode"];
            location?: string;
            message: string;
            metadata?: components["schemas"]["ErrorMetadata"];
        };
        ErrorMetadata: {
            allowed_values?: string[] | null;
            /** Format: int64 */
            max_length?: number;
            /** Format: double */
            max_value?: number;
            /** Format: int64 */
            min_length?: number;
            /** Format: double */
            min_value?: number;
            regex?: string;
        };
        ForgotPasswordRequestBody: {
            /** @description The Cloudflare Turnstile token */
            cfTurnstileToken: string;
            /**
             * Format: email
             * @description The user's email address
             */
            email: string;
        };
        GetAllSessionsResponseBody: {
            /** @description The current active user sessions */
            sessions: components["schemas"]["Session"][] | null;
        };
        MeResponseBody: {
            /** @description The currently active entity */
            activeEntity?: components["schemas"]["Entity"];
            /** @description The permissions for the currently active entity */
            entityRole?: components["schemas"]["EntityRole"];
            /** @description Whether two-factor authentication is pending */
            isTwoFactorPending: boolean;
            /** @description The user object */
            user: components["schemas"]["SessionUser"];
        };
        Membership: {
            /** @description The entity details */
            entity: components["schemas"]["Entity"];
            /** @description The entity ID */
            entityId: string | null;
            /** @description The membership ID */
            id: string;
            /** @description The user's role in the entity */
            role: string;
        };
        Payment: {
            /** Format: int64 */
            amount: number;
            /** Format: date-time */
            completed_at?: string;
            /** Format: date-time */
            created_at: string;
            currency: string;
            description: string;
            id: string;
            merchant_id: string;
            metadata: {
                [key: string]: unknown;
            };
            method: string;
            provider: string;
            status: string;
            /** Format: date-time */
            updated_at: string;
        };
        RegenerateQRCodeResponseBody: {
            /** @description The QR code for scanning with authenticator apps */
            qr_code: string;
        };
        ResetPasswordRequestBody: {
            /** @description The new password */
            newPassword: string;
            /**
             * Format: uuid
             * @description The password reset token
             */
            token: string;
        };
        Session: {
            /** @description The last seen IP country */
            country: string | null;
            /**
             * Format: date-time
             * @description The session creation time
             */
            createdAt: string;
            /** @description Whether this is the user's current session */
            current: boolean;
            /** @description The session ID */
            id: string;
            /** @description The last seen IP address */
            ipAddress: string | null;
            /**
             * Format: date-time
             * @description The session's last activity time
             */
            updatedAt: string;
            /** @description The last seen user agent */
            userAgent: string | null;
            /** @description The user's ID */
            userId: string;
        };
        SessionUser: {
            /** @description The user's email address */
            email: string;
            /** @description The user's ID */
            id: string;
            /** @description The user's profile image URL */
            image?: string;
            /** @description Whether two-factor authentication is enabled */
            isTwoFactorEnabled: boolean;
            /** @description Whether the user's email is verified */
            isVerified: boolean;
            /**
             * Format: date-time
             * @description The user's last activity time
             */
            lastActiveAt?: string;
            /**
             * Format: date-time
             * @description The user's last login time
             */
            lastLoggedInAt?: string;
            /** @description The user's entity memberships */
            memberships?: components["schemas"]["Membership"][] | null;
            /** @description The user's name */
            name: string;
            /**
             * Format: date-time
             * @description When the current session will expire
             */
            sessionExpiresAt: string;
        };
        SetupTwoFactorResponseBody: {
            /** @description The backup codes for account recovery */
            backupCodes: string[] | null;
            /** @description The QR code for scanning with authenticator apps */
            qrCode: string;
            /** @description The TOTP secret key */
            secret: string;
        };
        SignInRequestBody: {
            /** @description The Cloudflare Turnstile token */
            cfTurnstileToken: string;
            /**
             * Format: email
             * @description The user's email address
             */
            email: string;
            /** @description The user's password */
            password: string;
        };
        SignInResponseBody: {
            /** @description Whether two-factor authentication is pending */
            isTwoFactorPending: boolean;
        };
        SignUpRequestBody: {
            /** @description The Cloudflare Turnstile token */
            cfTurnstileToken: string;
            /**
             * Format: email
             * @description The user's email address
             */
            email: string;
            /** @description The user's full name */
            name: string;
            /** @description The user's password */
            password: string;
        };
        UpdatePasswordRequestBody: {
            /** @description The current password */
            currentPassword: string;
            /** @description The new password */
            newPassword: string;
        };
        UpdateUserRequestBody: {
            /** @description The user's full name */
            name: string;
        };
        User: {
            /** Format: date-time */
            createdAt: string;
            email?: string;
            /** Format: date-time */
            emailVerifiedAt?: string;
            id?: string;
            image?: string;
            /** Format: date-time */
            lastActiveAt?: string;
            /** Format: date-time */
            lastLoggedInAt?: string;
            name?: string;
            /** Format: date-time */
            updatedAt: string;
        };
        VerifyEmailRequestBody: {
            /**
             * Format: uuid
             * @description The verification token
             */
            token: string;
        };
        VerifyPasswordRequestBody: {
            /** @description The current password to verify */
            password: string;
        };
        VerifyPasswordResponseBody: {
            /** @description Whether the password was verified successfully */
            verified: boolean;
        };
        VerifyTwoFactorRequestBody: {
            /** @description The two-factor authentication code */
            code: string;
        };
    };
    responses: {
        /** @description Too many requests - rate limit exceeded */
        TooManyRequests: {
            headers: {
                /** @description The number of allowed requests in the current period */
                "X-RateLimit-Limit"?: number;
                /** @description The number of remaining requests in the current period */
                "X-RateLimit-Remaining"?: number;
                /** @description The remaining window before the rate limit resets in UTC epoch seconds */
                "X-RateLimit-Reset"?: number;
                [name: string]: unknown;
            };
            content: {
                "application/json": {
                    /** @description Error message */
                    error?: string;
                };
            };
        };
    };
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export interface operations {
    "disable-two-factor": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["DisableTwoFactorResponseBody"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "enable-two-factor": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EnableTwoFactorRequestBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["EnableTwoFactorResponseBody"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "forgot-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ForgotPasswordRequestBody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "get-session": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: {
                /** @description The session cookie */
                session?: string;
            };
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["MeResponseBody"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "refresh-session": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: {
                /** @description The refresh token cookie */
                refresh_token?: string;
            };
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    Domain?: string;
                    Expires?: string;
                    HttpOnly?: boolean;
                    MaxAge?: number;
                    Name?: string;
                    Partitioned?: boolean;
                    Path?: string;
                    Quoted?: boolean;
                    Raw?: string;
                    RawExpires?: string;
                    SameSite?: number;
                    Secure?: boolean;
                    "Set-Cookie"?: string;
                    Unparsed?: string;
                    Value?: string;
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "regenerate-qr-code": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["RegenerateQRCodeResponseBody"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "reset-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ResetPasswordRequestBody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "get-all-sessions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: {
                /** @description The session cookie */
                session?: string;
            };
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["GetAllSessionsResponseBody"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "delete-all-sessions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: {
                /** @description The session cookie */
                session?: string;
            };
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "delete-session": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description The session id */
                id: string;
            };
            cookie?: {
                /** @description The session cookie */
                session?: string;
            };
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "setup-two-factor": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SetupTwoFactorResponseBody"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "sign-in": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["SignInRequestBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    Domain?: string;
                    Expires?: string;
                    HttpOnly?: boolean;
                    MaxAge?: number;
                    Name?: string;
                    Partitioned?: boolean;
                    Path?: string;
                    Quoted?: boolean;
                    Raw?: string;
                    RawExpires?: string;
                    SameSite?: number;
                    Secure?: boolean;
                    "Set-Cookie"?: string;
                    Unparsed?: string;
                    Value?: string;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SignInResponseBody"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "sign-out": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: {
                /** @description The session cookie */
                session?: string;
            };
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    Domain?: string;
                    Expires?: string;
                    HttpOnly?: boolean;
                    MaxAge?: number;
                    Name?: string;
                    Partitioned?: boolean;
                    Path?: string;
                    Quoted?: boolean;
                    Raw?: string;
                    RawExpires?: string;
                    SameSite?: number;
                    Secure?: boolean;
                    "Set-Cookie"?: string;
                    Unparsed?: string;
                    Value?: string;
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "sign-up": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["SignUpRequestBody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "update-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpdatePasswordRequestBody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "verify-email": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["VerifyEmailRequestBody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "verify-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["VerifyPasswordRequestBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["VerifyPasswordResponseBody"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "verify-two-factor": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: {
                /** @description The session cookie */
                session?: string;
            };
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["VerifyTwoFactorRequestBody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    Domain?: string;
                    Expires?: string;
                    HttpOnly?: boolean;
                    MaxAge?: number;
                    Name?: string;
                    Partitioned?: boolean;
                    Path?: string;
                    Quoted?: boolean;
                    Raw?: string;
                    RawExpires?: string;
                    SameSite?: number;
                    Secure?: boolean;
                    "Set-Cookie"?: string;
                    Unparsed?: string;
                    Value?: string;
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "create-payment": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreatePaymentRequestBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Payment"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "get-user": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description The ID of the user. Use @me to refer to the current user.
                 * @example @me
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["User"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "update-user": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description The ID of the user. Use @me to refer to the current user.
                 * @example @me
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpdateUserRequestBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["User"];
                };
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
    "update-user-image": {
        parameters: {
            query?: never;
            header: {
                /** @description The original file name. */
                "X-File-Name": string;
            };
            path: {
                /**
                 * @description The ID of the user. Use @me to refer to the current user.
                 * @example @me
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "image/*": string;
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Too many requests - rate limit exceeded */
            429: components["responses"]["TooManyRequests"];
        };
    };
}
