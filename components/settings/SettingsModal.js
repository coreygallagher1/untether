import { Modal, Text, Button, Tooltip, Grid, Input } from "@nextui-org/react";
import { EditIcon } from "./EditIcon";
import { IconButton } from "./IconButton";

export default function SettingsModal(props) {
  return (
    <>
      <Modal
        closeButton
        aria-labelledby="modal-title"
        open={props.visible}
        onClose={props.closeHandler}
      >
        {/* <Modal.Header>
          <Text id="modal-title" size={18}>
            Settings
            <Text b size={18}>
              NextUI
            </Text>
          </Text>
        </Modal.Header> */}
        <Modal.Body>
          <Grid.Container gap={2}>
            <Grid xs={11}>
              <Input
                readOnly
                placeholder="Read only"
                initialValue="jacob@bitbybite.xyz"
                width="95%"
              />
            </Grid>
            <Grid xs={1}>
              <Tooltip content="Edit user">
                <IconButton>
                  <EditIcon size={20} fill="#979797" />
                </IconButton>
              </Tooltip>
            </Grid>
            <Grid xs={11}>
              <Input
                readOnly
                placeholder="Read only"
                initialValue="password"
                width="95%"
              />
            </Grid>
            <Grid xs={1}>
              <Tooltip content="Edit user">
                <IconButton>
                  <EditIcon size={20} fill="#979797" />
                </IconButton>
              </Tooltip>
            </Grid>
          </Grid.Container>
        </Modal.Body>
        {/* <Modal.Footer>
          <Button auto flat color="error" onPress={props.closeHandler}>
            Close
          </Button>
          <Button auto onPress={props.closeHandler}>
            Sign in
          </Button>
        </Modal.Footer> */}
      </Modal>
    </>
  );
}
