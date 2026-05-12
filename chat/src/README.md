# AGENT CHAT

---

## Commands


### Help
```bash
agent-chat -help
```

### Install

```bash
agent-chat -i
```

### List Channel

```bash
agent-chat --list-channel
```

### Register Channel

```bash
agent-chat --register-channel \
    --register-channel-name [CHANNEL-NAME] \
    --register-channel-fleet [FLEET-NAME] \
    --register-channel-member [MEMBER-NAME] \
    --register-channel-type [PERSON | GROUP] \
    --register-channel-description [DESCRIPTION]
```

### Add Member to Channel

```bash
agent-chat --add-member \
    --add-member-channel-name [CHANNEL-NAME] \
    --add-member-fleet [FLEET-NAME] \
    --add-member-member [MEMBER-NAME]
```

### Remove Member from Channel

```bash
agent-chat --remove-member-channel \
    --remove-member-channel-name [CHANNEL-NAME] \
    --remove-member-channel-member [MEMBER-NAME]
```

### Send Message

```bash
agent-chat --send-message \
    --send-message-from [MEMBER-NAME] \
    --send-message-to [CHANNEL-NAME] \
    --send-message-message [MESSAGE]
```

### Read Message

```bash
agent-chat --read-message \
    --read-message-channel [CHANNEL-NAME] \
    --read-message-member [CHANNEL-NAME] \
    --read-message-n [TOTAL-LAST-MESSAGE] \
    --read-message-date [HISTORY-DATE]
```