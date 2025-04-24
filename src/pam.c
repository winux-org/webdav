// +build linux

// it works with file ending file_linux.go (and c probably as well)

#define _POSIX_C_SOURCE 200809L  // Ensure strdup() is available

#include <security/pam_appl.h>
#include <stdlib.h>
#include <string.h>

// PAM Conversation Function
int pam_conv_func(int num_msg, const struct pam_message **msg,
                  struct pam_response **resp, void *appdata_ptr) {
    struct pam_response *responses = (struct pam_response *)malloc(num_msg * sizeof(struct pam_response));
    if (!responses) return PAM_CONV_ERR;

    char *password = (char *)appdata_ptr;

    for (int i = 0; i < num_msg; i++) {
        responses[i].resp_retcode = 0;
        if (msg[i]->msg_style == PAM_PROMPT_ECHO_OFF) {
            responses[i].resp = strdup(password);
        } else {
            responses[i].resp = NULL;
        }
    }

    *resp = responses;
    return PAM_SUCCESS;
}

// Function to authenticate user with PAM
int authenticate_pam(const char *username, const char *password) {
    struct pam_conv conv = { pam_conv_func, (void *)password };
    pam_handle_t *pamh = NULL;

    int retval = pam_start("login", username, &conv, &pamh);
    if (retval != PAM_SUCCESS) return retval;

    retval = pam_authenticate(pamh, 0);
    pam_end(pamh, retval);

    return retval;
}
