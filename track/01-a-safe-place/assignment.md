---
slug: a-safe-place
id: weeahhs2pihp
type: challenge
title: "\U0001F9BA This is your sandbox, a safe place"
teaser: You can't break anything
notes:
- type: text
  contents: |-
    In this sandbox you can use AWS with AdministratorAccess.

    The only catch is your subscription will be destroyed in 8hrs.
    You will be warned after 90mins of inactivity, so keep an eye out for that.

    So be sure to save any work outside of the sandbox, like to a version control system.
tabs:
- id: ajsgopsb4zth
  title: VerteX
  type: service
  hostname: vertex-client
  path: /
  port: 80
- id: imodedbnyuvw
  title: Shell
  type: terminal
  hostname: vertex-client
difficulty: basic
enhanced_loading: null
---

<h2>Create a key pair</h2>


Before using Palette to deploy into AWS
we need to create a key pair that will be injected into any of the ec2 resources.

<h2>Let's use the AWS CLI</h2>

The AWS CLI will allow us to create a key pair without having to make a bunch of clicks.

- In your Shell tab run the following commands

- The first command will set our AWS operating region, in this case us-east-1

- The second command will create a key pair in us-east-1, named **MyKeyPair**
***(feel free to change the name)***

```bash
aws configure set default.region us-east-1
```

**If you just want to use these creds for a UI demo create a key pair first**

```bash
aws ec2 create-key-pair --key-name MyKeyPair
```
After the key generates

Press `q`

to regain terminal control

<h2>🛠 CLI Tools</h2>


This shell has the following CLI Tools installed

- AWS CLI v2
- kubectl
- Helm
- Terraform

<h2>Terraform</h2>


To use Terraform in this sandbox you can clone a repo
into the /tmp directory using the Shell tab

<h2>Screen</h2>


The Shell tab has Screen installed
You can create more terminal sessions using

`ctrl+a+c`

You can switch between terminal session using

`shift+right arrow or left arrow`

