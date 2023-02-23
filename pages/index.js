import { Image, Text, Grid, Button } from "@nextui-org/react";
import { useRouter } from "next/router";


export const AcmeLogo = () => (
  <svg
    className=""
    fill="none"
    height="36"
    viewBox="0 0 32 32"
    width="36"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect fill="var(--secondary)" height="100%" rx="16" width="100%" />
    <path
      clipRule="evenodd"
      d="M17.6482 10.1305L15.8785 7.02583L7.02979 22.5499H10.5278L17.6482 10.1305ZM19.8798 14.0457L18.11 17.1983L19.394 19.4511H16.8453L15.1056 22.5499H24.7272L19.8798 14.0457Z"
      fill="currentColor"
      fillRule="evenodd"
    />
  </svg>
);

export default function Index() {
  const router = useRouter();
  return (
    <>
      <Grid.Container gap={2} justify="center">
        <Grid xs={50} justify="center">
          <Text style={{ fontSize: 50, color: '#6D6C95', alignContent: "center", 
          weight: "bold" , justify: "flex-end"}} 
           h2>Do you want to get involved in Milwaukee's Tech hub?
           </Text>
        </Grid>
        <Grid xs={12} style={{ marginBottom: 30, marginLeft: 75, marginRight: 75}}>
          <text style={{ fontSize: 20, color: '#6D6C95', justify: "center", align :"center" }}>
          Our mission at The Milky Way Tech Hub is to support the growth of Milwaukee's 
          tech community. We give young entrepreneurs the support and resources they need to help 
          expand their businesses and participate in innovation. It is important to us that we create equitable opportunities 
          and inspire diverse individuals to break into tech.

          In support of our mission, we have decided to make this Ecosytem app free to provide you with resources to help you learn more about
          the tech space, how to achieve your business goals, and more. 
          </text>
        </Grid>
        <img width={1000} height={500} src="https://i0.wp.com/milkywaytechhub.com/wp-content/uploads/2019/02/fly-me-to-the-moon-696x392.png" 
        style={{ marginBottom: 30, marginLeft: 75, marginRight: 75}}/>
 
        <Grid xs={4} justify="center">
          <Text style={{ fontSize: 20}} className='Purple' h3>Looking for Resources?</Text>
        </Grid>
        <Grid xs={4} justify="center">
          <Text style={{ fontSize: 20}} className='Purple' h3>Get Custom Business Guides</Text>
        </Grid>
        <Grid xs={4} justify="center">
          <Text style={{ fontSize: 20}} className='Purple' h3>Expand your Network</Text>
        </Grid>
        <Grid xs={4} justify="center">
          <Button css={{ background: '#6D6C95'}} onClick={() => router.push("/resources")}>
            Resources
          </Button>
        </Grid>
        <Grid xs={4} justify="center">
          <Button css={{ background: '#6D6C95'}} onClick={() => router.push("/Quiz")}> Take Questionare</Button>
        </Grid>
        <Grid xs={4} justify="center">
          <Button css={{ background: '#6D6C95'}} onClick={() => router.push("https://milkywaytechhub.com/groups/")}>Connect</Button>
       
        </Grid>
      </Grid.Container>
    </>
  );
}
