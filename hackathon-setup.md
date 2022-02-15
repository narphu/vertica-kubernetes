Here are steps to get the MC deployment working for the hackathon:

1. Install go.  Must be version 1.16.2.

curl -LO https://go.dev/dl/go1.16.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz

2. Install docker

3. Install kubectl

curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

4. Install helm

curl -LO https://get.helm.sh/helm-v3.8.0-linux-amd64.tar.gz
tar zxvf helm-v3.8.0-linux-amd64.tar.gz
sudo install -o root -g root -m 0755 linux-amd64/helm /usr/local/bin/helm

5. Install kubectx

curl -LO https://github.com/ahmetb/kubectx/releases/download/v0.9.4/kubectx_v0.9.4_linux_x86_64.tar.gz
tar zxvf kubectx_v0.9.4_linux_x86_64.tar.gz
sudo install -o root -g root -m 0755 kubectx /usr/local/bin/kubectx

6. Install kubens

curl -LO https://github.com/ahmetb/kubectx/releases/download/v0.9.4/kubens_v0.9.4_linux_x86_64.tar.gz
tar zxvf kubens_v0.9.4_linux_x86_64.tar.gz
sudo install -o root -g root -m 0755 kubens /usr/local/bin/kubens

7. Install krew

(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
)

8. Add $HOME/.krew/bin/ to your path in your .bashrc (or whatever login shell you use)

9. Restart shell

10. Change directory in your shell to this git repository

11a. Remove an old kind cluster (if these steps were done before).

scripts.kind.sh term cluster1

11b. Start kubernetes-in-docker (kind).

scripts/kind.sh init cluster1  

If you intend to use nodePort, use the the following command instead:

scripts/kind.sh -p 30001 init cluster1  

12. Setup minio in k8s

scripts/setup-minio.sh -o

13. Install cert-manager in k8s

make install-cert-manager

14. Copy console rpm to docker-console/package/vertica-console-x86_64.RHEL6.latest.rpm

15. Build the operator and console containers

make docker-build-operator docker-build-console docker-push-operator docker-push-console

16. Deploy the operator and the console

make deploy

17. Pre-push the vertica server image

scripts/push-to-kind.sh -i vertica/vertica-k8s:latest

18. Create a minIO tenants.  This will serve as the communal path.

kubectl apply -f config/samples/minio.yaml

19. Wait for minIO to come up.  Run this in a loop until you see one of the pods have STATUS column set to completed.

kubectl get pods --selector job-name=create-s3-bucket
NAME                     READY   STATUS      RESTARTS   AGE
create-s3-bucket-4qprk   0/1     Error       0          42s
create-s3-bucket-bbtpp   0/1     Error       0          52s
create-s3-bucket-cvz5q   0/1     Error       0          104s
create-s3-bucket-jppz6   0/1     Completed   0          22s

20. Create the Vertica DB

kubectl apply -f config/samples/v1beta1_verticadb.yaml

21.  Wait for the database to be created

kubectl wait --for=condition=DBInitialized=true vdb/verticadb-sample --timeout 600s


If you are going to use nodePort skip the next 2 steps.

22. Expose the MC so you can access it on your localhost

scripts/expose-console.sh

23.  Open up MC in your webbrowser

https://localhost:5450/

If you are going to use nodePort do these two steps instead of the last two.

24.  Update Service object to be NodePort

kubectl patch svc verticadb-operator-console --type='json' -p '[{"op":"replace","path":"/spec/type","value":"NodePort"},{"op":"replace","path":"/spec/ports/0/nodePort","value":30001}]]'

25. Open up MC in your webbrowser.  Pay attention to the port number it is different than in step 23.  5433 is the port on your hosts that maps to 30001 (in docker) that was setup when creating the kind cluster.

https://localhost:5433/
