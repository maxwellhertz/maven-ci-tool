package app_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/maxwellhertz/maven-project-ci-tool/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockMavenRepoApi struct {
	mock.Mock
}

func (m *mockMavenRepoApi) GetLatestReleaseVersion(artifactId string) (string, error) {
	args := m.Called(artifactId)
	return args.String(0), args.Error(1)
}

const validInitialPom = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
	<modelVersion>4.0.0</modelVersion>
</project>`

func TestConfigurePluginsSucceed(t *testing.T) {
	testCases := []struct {
		desc string
		xml  string
	}{
		{
			desc: "missing <build> element",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
			<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
				<modelVersion>4.0.0</modelVersion>
				<groupId>com.example</groupId>
				<artifactId>demo</artifactId>
				<version>0.0.1-SNAPSHOT</version>
				<name>demo</name>
			</project>
			`,
		},
		{
			desc: "missing <plugins> element",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
			<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
				<modelVersion>4.0.0</modelVersion>
				<groupId>com.example</groupId>
				<artifactId>demo</artifactId>
				<version>0.0.1-SNAPSHOT</version>
				<name>demo</name>
				<build></build>
			</project>`,
		},
		{
			desc: "missing maven-surefire-plugin",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
			<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
				<modelVersion>4.0.0</modelVersion>
				<groupId>com.example</groupId>
				<artifactId>demo</artifactId>
				<version>0.0.1-SNAPSHOT</version>
				<name>demo</name>
				<build>
					<plugins></plugins>
				</build>
			</project>`,
		},
		{
			desc: "maven-surefire-plugin exists",
			xml: `<?xml version="1.0" encoding="UTF-8"?>
			<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
				<modelVersion>4.0.0</modelVersion>
				<groupId>com.example</groupId>
				<artifactId>demo</artifactId>
				<version>0.0.1-SNAPSHOT</version>
				<name>demo</name>
				<build>
					<plugins>
						<plugin>
							<groupId>org.apache.maven.plugins</groupId>
							<artifactId>maven-surefire-plugin</artifactId>
							<version>0.1.0-SNAPSHOT</version>
						</plugin>
					</plugins>
				</build>
			</project>`,
		},
	}

	mockApi := new(mockMavenRepoApi)
	mockApi.On("GetLatestReleaseVersion", mock.Anything).Return("1.0.0", nil)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			u := app.NewDefaultPomUitls(mockApi)
			_, err := u.ConfigurePlugins(strings.NewReader(tc.xml))
			assert.Nil(t, err)
			mockApi.AssertExpectations(t)
		})
	}
}

func TestMavenRepositoryApiFail(t *testing.T) {
	mockApi := new(mockMavenRepoApi)
	mockApi.On("GetLatestReleaseVersion", mock.Anything).Return(mock.Anything, errors.New("some error"))

	u := app.NewDefaultPomUitls(mockApi)
	_, err := u.ConfigurePlugins(strings.NewReader(validInitialPom))
	assert.NotNil(t, err)
	mockApi.AssertExpectations(t)
}

func TestParseInvliadPomFilet(t *testing.T) {
	mockApi := new(mockMavenRepoApi)

	testCases := []struct {
		xml  string
		desc string
	}{
		{
			xml:  `<?xml version="1.0" encoding="UTF-8"?>`,
			desc: "missing project",
		},
		{
			xml: `<?xml version="1.0" encoding="UTF-8"?>
			<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
			</project>`,
			desc: "missing modelVersion",
		},
		{
			xml: `<?xml version="1.0" encoding="UTF-8"?>
			<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
				<modelVersion>1.0.0</modelVersion>
			</project>`,
			desc: "invalid modelVersion",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			u := app.NewDefaultPomUitls(mockApi)
			_, err := u.ConfigurePlugins(strings.NewReader(tc.xml))
			assert.NotNil(t, err)
		})
	}
}