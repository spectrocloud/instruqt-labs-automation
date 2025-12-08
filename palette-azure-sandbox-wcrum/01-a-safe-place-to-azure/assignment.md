---
slug: a-safe-place-to-azure
id: yyabchulwtqx
type: challenge
title: "\U0001F9BA This is your Azure sandbox, a safe place"
teaser: You can't break anything
notes:
- type: text
  contents: |-
    In this sandbox you can use an Azure subscription with Spectro Cloud Palette

    The only catch is your subscription will be destroyed in 24hrs

    So be sure to save any work outside of the sandbox, like to a version control system
tabs:
- id: yufxjemyeern
  title: Palette
  type: service
  hostname: palette-client
  path: /
  port: 80
- id: 9xsmovn06l2w
  title: Cloud CLI
  type: terminal
  hostname: cloud-client
- id: mmsx399h0rxy
  title: Azure Portal
  type: service
  hostname: cloud-client
  path: /
  port: 80
- id: xybqbpv3mpyp
  title: Text Editor
  type: code
  hostname: cloud-client
  path: /root
- id: 73r8ax2jyp9w
  title: Spectro Cloud Palette Trial
  type: website
  url: https://www.spectrocloud.com/free-trial/
  new_window: true
difficulty: basic
enhanced_loading: null
---

<h2>Prepping your environment</h2>

Before using Palette to deploy into the Azure subscription
we need to create an Azure Resource Group that will be used for creating K8s infrastructure

<h2>Using the Azure CLI</h2>


The Azure CLI is pre-installed for you and logged in with the temp credentials provided

In the terminal run this command

```bash
az group create --name MyResourceGroup --location eastus
```

You should see this type of output

```
{
  "id": "/subscriptions/5629542c-39a2-4bce-81bf-cbebf4f1c511/resourceGroups/MyResourceGroup",
  "location": "eastus",
  "managedBy": null,
  "name": "MyResourceGroup",
  "properties": {
    "provisioningState": "Succeeded"
  },
  "tags": null,
  "type": "Microsoft.Resources/resourceGroups"
}
```

*you can customize the --location to an Azure region of your choice

Good job creating a Resource Group, let's move on.




<h2>My Azure CLI is broken ☹️</h2>


If you were not able to create the resource group your Azure cli may not be logged in

Try running this command

```bash
az login --service-principal -u $ARM_CLIENT_ID -p=$ARM_CLIENT_SECRET --tenant $ARM_TENANT_ID
```

Then try again

```bash
az group create --name MyResourceGroup --location eastus
```

<h2>Configure a cloud account in Palette</h2>

Create a Spectro Cloud Palette trial or log into your existing trial account

Navigate to the Project Settings> Cloud Accounts > Add Azure Account

![Add Azure Account](../assets/palette_cloud-account.png)

In the Instruqt tab **Azure Portal** copy your cloud credentials

and paste them into the Palette configuration wizard

Follow this image for the corresponding values

![Palette Azure Account setup](../assets/palette-cloud-account-fields.png)

Click **Validate** then **Confirm**