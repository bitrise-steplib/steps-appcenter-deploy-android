# AppCenter Android Deploy

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-appcenter-deploy-android?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-appcenter-deploy-android/releases)

Distribute your Android app through [Microsoft App Center](https://appcenter.ms/).

<details>
<summary>Description</summary>

[This Step](https://github.com/bitrise-steplib/steps-appcenter-deploy-android) integrates with the [App Center](https://appcenter.ms/)'s Distribution service and enables you to distribute your apps seamlessly to different stores, for example, App Store, MS Intune, user groups or even individual testers.

 ### Configuring the Step

 Before you start:

 The Step requires an active MS App Center account.

 1. Add the **APP path** which points to a binary file.
 2. Add the **mapping.txt file path**.
 3. Add the App Center **API token**.
 4. Add the **Owner name**, which means the owner of the App Center app. For an app owned by a user, the URL in App Center can look like this https://appcenter.ms/users/JoshuaWeber/apps/APIExample where the {ownername} is JoshuaWeber. For an app owned by an organization, the URL can be, for example, https://appcenter.ms/orgs/Microsoft/apps/APIExample where the {ownername} is Microsoft.
 5. Add the **App name** which is the name of the App Center app. For an app owned by a user, the URL in App Center might look like this: https://appcenter.ms/users/JoshuaWeber/apps/APIExample where the {app_name} is APIExample.
 6. Add the **Distribution groups** which means the user groups you wish to distribute the app to. Please add one group name per line.
 7. Add the **Distribution stores** where you wish to distribute the app to. Please add one store name per line.
 8. Add the **Testers** who you wish to send the app to via email. Please add one email address per line.
 9. Add any **Release notes for the deployed artifact**.
 10. Send notification emails to testers and distribution groups with the **Notify Testers** input.
 11. You can enforce the installation of a distribution version with the **Mandatory** input set to `yes`.
 12. If you set the **Debug** input to `yes`, you can enable verbose logs.

### Useful links
- [About Android deployment with Bitrise](https://devcenter.bitrise.io/deploy/android-deploy/android-deployment-index/)
- [About Android code signing](https://devcenter.bitrise.io/code-signing/android-code-signing/android-code-signing-index/)

### Related Steps
- [Deploy to Huawei App Gallery](https://www.bitrise.io/integrations/steps/app-gallery-deploy)
- [Google Play Deploy](https://www.bitrise.io/integrations/steps/google-play-deploy)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `app_path` | Path to binary file  For APKs, only single or universal APKs are supported: https://docs.microsoft.com/en-us/appcenter/build/react-native/android/#63-building-multiple-apks | required | `$BITRISE_APP_PATH` |
| `mapping_path` | Path to an Android mapping.txt file. |  |  |
| `api_token` | App Center API token | required, sensitive |  |
| `owner_name` | Owner of the App Center app.  For an app owned by a user, the URL in App Center might look like https://appcenter.ms/users/JoshuaWeber/apps/APIExample.  Here, the {owner_name} is JoshuaWeber. For an app owned by an org, the URL might be https://appcenter.ms/orgs/Microsoft/apps/APIExample and the {owner_name} would be Microsoft | required |  |
| `app_name` | The name of the App Center app.  For an app owned by a user, the URL in App Center might look like https://appcenter.ms/users/JoshuaWeber/apps/APIExample.  Here, the {app_name} is ApiExample. | required |  |
| `distribution_group` | User groups you wish to distribute the app. One group name per line.  Distribution of AAB is supported only for Google Play store deployment: https://docs.microsoft.com/en-us/appcenter/distribution/uploading#android |  |  |
| `distribution_store` | Distribution stores you wish to distribute the app. One store name per line.  Distribution of AAB is supported only for Google Play store deployment: https://docs.microsoft.com/en-us/appcenter/distribution/uploading#android |  |  |
| `distribution_tester` | List of individual testers. One email per line.  Distribution of AAB is supported only for Google Play store deployment: https://docs.microsoft.com/en-us/appcenter/distribution/uploading#android |  |  |
| `release_notes` | Additional notes for the deployed artifact. |  | `Release notes` |
| `notify_testers` | Send notification email to testers and distribution groups. | required | `yes` |
| `mandatory` | Enforce installation of distribution version. Requires SDK integration. | required | `no` |
| `debug` | Enable verbose logs | required | `no` |
| `all_distribution_groups` | Distribute the app to all user groups on that app. Enabling this options makes it ignore distribution_group. |  | `no` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `APPCENTER_DEPLOY_STATUS` | Deployment status: 'success' or 'failed' |
| `APPCENTER_DEPLOY_INSTALL_URL` | Install page URL of the newly deployed version. |
| `APPCENTER_DEPLOY_DOWNLOAD_URL` | Download URL of the newly deployed version. |
| `APPCENTER_DEPLOY_RELEASE_ID` | ID of the new release for later retrieval via App Center APIs. |
| `APPCENTER_PUBLIC_INSTALL_PAGE_URL` | Public install page URL of the latest version. |
| `APPCENTER_PUBLIC_INSTALL_PAGE_URLS` | When a group is public the step will AppCenter provides and the step exports a public install page URL. |
| `APPCENTER_RELEASE_PAGE_URL` | URL to the release page containing release notes, easily share with business partners and QA for testing. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-appcenter-deploy-android/pulls) and [issues](https://github.com/bitrise-steplib/steps-appcenter-deploy-android/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

**Note:** this step's end-to-end tests (defined in `e2e/bitrise.yml`) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
