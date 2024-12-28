package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/iamYole/gsocial/internal/store"
)

var usernames = []string{
	"TechGuru01", "CodeWizard02", "DevMaster03", "GoLangNinja04", "DataCruncher05",
	"CloudHunter06", "CryptoKnight07", "BugSquasher08", "PixelPusher09", "AppArchitect10",
	"ServerSage11", "TerminalWizard12", "CyberExplorer13", "AlgoPro14", "StackOverflow15",
	"BinaryBoss16", "DebugHero17", "AI_Maestro18", "UX_Crafter19", "BackendBuff20",
	"FrontEndFreak21", "DevOpsDynamo22", "CodeSniper23", "HackyHacker24", "LogicLover25",
	"ScriptGuru26", "NullPointer27", "BitStreamer28", "CompileKing29", "LazyLoader30",
	"MemoryMiner31", "FuncFan32", "ThreadMaster33", "CacheLord34", "AsyncAce35",
	"QueueConqueror36", "HTTPHero37", "RegexRuler38", "LoopWizard39", "BuildBreaker40",
	"PatchHunter41", "SyntaxSlayer42", "PackagePirate43", "CloudCoder44", "ModuleMaker45",
	"GoGetter46", "InterfaceImp47", "ConcurrentKid48", "PointerPal49", "ChannelChampion50",
}
var titles = []string{
	"Canary Releases: Minimizing Deployment Risks", "DevOps with AWS: Key Services and Use Cases",
	"DevOps Culture: Fostering Collaboration", "Creating Effective Postmortem Reports", "API Gateways and DevOps",
	"Advanced Kubernetes Autoscaling Techniques", "Using Jenkins Pipelines for Automation", "Managing Kubernetes with kustomize",
	"Building CI/CD Pipelines with GitHub Actions", "Implementing DevOps in Legacy Applications", "Site Reliability Engineering Tools and Techniques",
	"Deploying Stateful Applications on Kubernetes", "Using Elastic Stack (ELK) for DevOps Monitoring", "Container Security Best Practices",
	"Handling Alerts Effectively in DevOps", "Designing Effective Release Management Processes", "Configuration Drift Detection and Remediation",
	"Understanding Load Testing in DevOps", "Using HashiCorp Consul for Service Discovery", "Building Resilient Systems with Circuit Breakers",
	"Introduction to Docker Swarm", "Compliance Automation in DevOps", "Secrets Management in Kubernetes", "Managing Logs and Metrics in Hybrid Clouds",
	"DevOps for Edge Computing Environments", "Automating Cloud Cost Monitoring", "DevOps Strategies for Continuous Feedback", "Comparing Kubernetes Operators vs. Helm Charts",
	"Building Serverless Pipelines with AWS Lambda", "Leveraging AI in DevOps Workflows", "Container Orchestration with Nomad",
	"Building Multi-Stage CI/CD Pipelines", "Performance Tuning for CI/CD Pipelines", "DevOps Compliance for Highly Regulated Industries",
	"Building Effective Chaos Experiments", "Kubernetes Cluster Upgrades Best Practices", "Integrating DevOps and ITIL Frameworks", "Using Pulumi for Infrastructure as Code",
}
var contents = []string{
	"DevOps bridges the gap between development and operations teams. It ensures faster delivery and higher reliability of applications.",
	"CI/CD pipelines automate code integration and deployment. This speeds up the release cycle and minimizes human errors.",
	"Kubernetes orchestrates containerized applications. It simplifies scaling, deployment, and management in the cloud.",
	"Infrastructure as Code allows version-controlled infrastructure. Tools like Terraform and CloudFormation make it possible.",
	"Monitoring is crucial for maintaining system health. Tools like Prometheus and Grafana help visualize metrics effectively.",
	"DevSecOps integrates security into the DevOps process. It ensures vulnerabilities are caught early in the pipeline.",
	"Automation is at the heart of DevOps efficiency. It reduces manual effort and increases consistency in workflows.",
	"GitOps uses Git as a single source of truth for deployments. It simplifies managing infrastructure and application updates.",
	"Blue-Green deployments minimize downtime during releases. Traffic is routed between two environments to ensure availability.",
	"Service Mesh like Istio manages communication between microservices.It enhances security, observability, and traffic control.",
	"Containerization ensures consistent application environments. Docker is the most popular tool for building and managing containers.",
	"Load balancing ensures even distribution of traffic. This improves performance and prevents server overloads.",
	"Cloud-native applications are designed for scalability. They leverage cloud services and containerization technologies.",
	"Secrets management protects sensitive credentials and keys. Tools like HashiCorp Vault help securely store and retrieve them.",
	"Log aggregation consolidates logs for easier analysis. The ELK stack (Elasticsearch, Logstash, Kibana) is a popular choice.",
	"Site Reliability Engineering (SRE) focuses on reliability and uptime. It applies engineering principles to operations tasks.",
	"Canary deployments test new releases on a small audience first. This reduces the risk of rolling out faulty updates.",
	"Observability combines monitoring, logging, and tracing. It helps diagnose and resolve issues in distributed systems.",
	"Serverless computing removes the need to manage servers. Platforms like AWS Lambda allow running code without provisioning infrastructure.",
	"Feature flags enable toggling features on or off in real time.This allows gradual rollouts and easier debugging of features.",
}
var tags = []string{
	"DevOps", "CI/CD", "Kubernetes", "Docker", "InfrastructureAsCode", "Monitoring", "CloudComputing", "GitOps", "DevSecOps", "Microservices",
	"Terraform", "Ansible", "AWS", "Azure", "GoogleCloud", "Automation", "Observability", "SiteReliabilityEngineering", "Containers", "LoadBalancing",
	"SecretsManagement", "CloudNative", "Serverless", "BlueGreenDeployments", "CanaryReleases", "Logging", "Prometheus", "Grafana", "ChaosEngineering", "ServiceMesh",
}
var comments = []string{
	"Great job on automating the CI/CD pipeline!", "Consider improving the monitoring setup for better insights.",
	"The Kubernetes deployment looks solid.", "Can we add more logging for debugging purposes?", "The new feature rollout was seamless—well done!",
	"Let's ensure we have proper rollback strategies in place.", "Have we accounted for scaling under peak load conditions?",
	"Fantastic use of Terraform for managing infrastructure.", "Is the deployment pipeline secure from unauthorized access?",
	"Don't forget to document the changes in the wiki.", "Could we optimize the container image size further?",
	"The metrics dashboard looks very informative.", "Have we stress-tested the system with high traffic?",
	"Good idea to use blue-green deployments for this release.", "The Helm charts need better versioning practices.", "Make sure all secrets are stored securely in the vault.",
	"The service mesh configuration is working well.", "Can we add more tests for the new microservice?", "The observability setup has greatly improved troubleshooting.",
	"Ensure the servers are patched with the latest updates.", "Using Prometheus for monitoring was a smart choice.", "Have we checked for compliance with security policies?",
	"The downtime was minimal during the deployment—great work!", "Can we enable feature toggles for this functionality?", "The autoscaling configuration seems to be working perfectly.",
	"Don't forget to clean up unused resources in the cloud.", "Consider using canary releases for the next update.", "Ensure all team members are familiar with the new process.",
	"The chaos testing results were insightful—good initiative!", "Have we updated the runbooks to reflect the changes?",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
		}
	}

	posts := generatePosts(300, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating user: ", err)
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123456",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: comments[rand.Intn(len(comments))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}
func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comments {
	cms := make([]*store.Comments, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comments{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
