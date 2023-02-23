import { Text, Grid } from "@nextui-org/react";
import { useRouter } from "next/router";

export default function Resources() {
  const router = useRouter();
  return (
    <>
      <Grid.Container gap={2}>
        <Grid xs={12}>
          <Text h1>Resources</Text>
        </Grid>
        <Grid xs={12}>
          <Text h3>Tech Events</Text>
        </Grid>
        <Grid xs={12}>
          <Text h3>MKE TECH NEWS</Text>
        </Grid>
        <Grid xs={12}>
          <Text h3>Podcasts</Text>
        </Grid>
        <Grid xs={12}>
          <Text h3>Business First Steps</Text>
        </Grid>
        <Grid xs={12}>
          <Text h3>Angel/Investment Firms</Text>
        </Grid>
        <Grid xs={12}>
          <Text h3>Meet Other Local Businesses</Text>
        </Grid>
      </Grid.Container>
    </>
  );
}
