package steamgamedl

import (
	"io"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		os   string
		body io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Linux",
			args: args{
				os:   "linux",
				body: strings.NewReader("\"linux\"\n{\n\t\"version\"\t\t\"1669935972\"\n\t\"steamcmd_public_all\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_public_all.zip.b34b8da97a76159a60dc27b1b6da62b538602518\"\n\t\t\"size\"\t\t\"49904\"\n\t\t\"sha2\"\t\t\"5320e6f2b60ec552de0abae4cf81df6bdeb4f29742bfa2ecdeaa94b067cff163\"\n\t}\n\t\"steamcmd_bins_linux\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_bins_linux.zip.71d25225b9323defa632d86eedfe7612d0660509\"\n\t\t\"size\"\t\t\"27543075\"\n\t\t\"sha2\"\t\t\"2793f5ea43e13c042f878e0a45308df3be58e5e4026f76fe0d9694abd0bc3ed5\"\n\t\t\"zipvz\"\t\t\"steamcmd_bins_linux.zip.vz.942f2cc97627404c2cc170a46c8752afa92c67db_19153875\"\n\t\t\"sha2vz\"\t\t\"f26cb5306bc12033748608623f42e04e96fda9b548e615199a6bf760b6b74af0\"\n\t}\n\t\"steamcmd_siteserverui_linux\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_siteserverui_linux.zip.d8b74ed2b27cf573238f13d6c207cef30a93c846\"\n\t\t\"size\"\t\t\"56381418\"\n\t\t\"sha2\"\t\t\"d56a26eaed091aad5c9285d2b7ed77cb0f76dc1b0da6f47b8305e7e2e7cf59d5\"\n\t\t\"zipvz\"\t\t\"steamcmd_siteserverui_linux.zip.vz.fc634823dd1831a88610f41bea8167a369bebb12_37365774\"\n\t\t\"sha2vz\"\t\t\"b721e42e18c1df26b709cec3e06950fc343ace3f4eedf55fefb1d6f4d07cfaaa\"\n\t}\n\t\"steamcmd_linux\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_linux.zip.29a467199e0d25ab89377b605511bd913d9cc94b\"\n\t\t\"size\"\t\t\"2739291\"\n\t\t\"sha2\"\t\t\"242458cf56972c62ccef324899daefaa516b96d5b5d3f20bbe82d038bdea8c00\"\n\t\t\"zipvz\"\t\t\"steamcmd_linux.zip.vz.88980b3581a48c17c9230ec19ff71373ad530ad9_2193605\"\n\t\t\"sha2vz\"\t\t\"db1e0736128f3276d7eba46876028ec79e5e64e04e0f02fb0c6039c9e864679e\"\n\t\t\"IsBootstrapperPackage\"\t\t\"1\"\n\t}\n}\n\"kvsign2\"\n{\n\t\"linux\"\t\t\"4670134441aa1dd10512f913db64c631a9b0b88be25a18910a7bebe3aaeefa30f3c1011376fcdcab76eed14cfbf63f03789bb7c82b0e7ea98897144d5fbd6208\"\n}\n\"kvsignatures\"\n{\n\t\"linux\"\t\t\"aba7b21e7b09d2d4707532dcb1638024b8a325d9323c180e371a5bf981da61e6ec923674b305d54fbda2d8c516734d2ef7b598a706456fb4c014af63949083062ab07f1fcf0705a290ea1d3ad585b58d0bc27bd4a90df30602b7a580704e6e61d9a5d6bf4d4a804d53ce3adc091ddcdc3667fd766cf93fbec0dcff695c9afc69\"\n}\n"),
			},
			want:    "steamcmd_bins_linux.zip.71d25225b9323defa632d86eedfe7612d0660509",
			wantErr: false,
		},
		{
			name: "Windows",
			args: args{
				os:   "win32",
				body: strings.NewReader("\"win32\"\n{\n\t\"version\"\t\t\"1669935972\"\n\t\"steamcmd_public_all\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_public_all.zip.b34b8da97a76159a60dc27b1b6da62b538602518\"\n\t\t\"size\"\t\t\"49904\"\n\t\t\"sha2\"\t\t\"5320e6f2b60ec552de0abae4cf81df6bdeb4f29742bfa2ecdeaa94b067cff163\"\n\t}\n\t\"steamcmd_bins_win32\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_bins_win32.zip.45ece1457835ba765444814e10b251645e1a8827\"\n\t\t\"size\"\t\t\"10258459\"\n\t\t\"sha2\"\t\t\"08afba6459b89f8b25562c958d0ea9c65dc5f58fbebd0cc73d7f86f0b58025c5\"\n\t\t\"zipvz\"\t\t\"steamcmd_bins_win32.zip.vz.d6bdd64485789afb88116d2273246206d1b91230_6413694\"\n\t\t\"sha2vz\"\t\t\"3cf37b3e0b25845bbb5e9f1487fefe412eeb17dd158aab7fbcf58b0d529c4c48\"\n\t}\n\t\"steamcmd_steamservice_win32\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_steamservice_win32.zip.a7099792cf9ca3b8731cf1163515588aa5221d21\"\n\t\t\"size\"\t\t\"2984037\"\n\t\t\"sha2\"\t\t\"41f01f06fe4eabb63d39eb71226e16991684dec9a1ac944e7d7d63e02cd599cb\"\n\t\t\"zipvz\"\t\t\"steamcmd_steamservice_win32.zip.vz.8d6d3753837c58e9885f1b154f5c8b2e1ad5569d_1682854\"\n\t\t\"sha2vz\"\t\t\"a639c2058a4ed404a7c821c588b9370bce981e28ca3ed2c8f03c7c05a5699d83\"\n\t}\n\t\"steamcmd_steamerrorreporter_win32\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_steamerrorreporter_win32.zip.9116eb26294f5997dea60a79117025c08668b6b9\"\n\t\t\"size\"\t\t\"184761\"\n\t\t\"sha2\"\t\t\"744ed1a40f0c954845550c5e5ea070bb56fde2973c2efc5c7225994497435049\"\n\t}\n\t\"steamcmd_siteserverui_win32\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_siteserverui_win32.zip.0dab448a609133fe0c361eb5afbbfc7fedf77602\"\n\t\t\"size\"\t\t\"51127464\"\n\t\t\"sha2\"\t\t\"8dedaf6928d0c3f66c6bc79bf5ebf8cedf3e7c86220556449f46805ba6ea7cc4\"\n\t\t\"zipvz\"\t\t\"steamcmd_siteserverui_win32.zip.vz.71927eb8e3fb154268d3dbf0014214f1fbc4966f_34465489\"\n\t\t\"sha2vz\"\t\t\"3d061c4fe911b781c640105875e4342718ad840cd6b5b365519a985e352d0408\"\n\t}\n\t\"steamcmd_win32\"\n\t{\n\t\t\"file\"\t\t\"steamcmd_win32.zip.bf29107b6df31bfa3962bc285033eef29b1c0212\"\n\t\t\"size\"\t\t\"2057969\"\n\t\t\"sha2\"\t\t\"f00352213d5aad53a0975724fc75605a4afb959d8b2b6c692942a95a903bd98a\"\n\t\t\"zipvz\"\t\t\"steamcmd_win32.zip.vz.c7d083accbea956356295383755feb70b3763657_1690813\"\n\t\t\"sha2vz\"\t\t\"0d99e7e828412e76b4b3b1d721d5cd4ee34fbd16cc10602556f2e4f119c4e2d5\"\n\t\t\"IsBootstrapperPackage\"\t\t\"1\"\n\t}\n}\n\"kvsign2\"\n{\n\t\"win32\"\t\t\"6aa237eac034527389e68eaad2cd21793671a0c934fc681483eabd3b98e8ea743e563c71ae62abf00f107dc70f23dd1b9c4c416b9e0a0d7ec6a861fa82998c01\"\n}\n\"kvsignatures\"\n{\n\t\"win32\"\t\t\"16a71a9d316a9273912a412ed9c4840f48be2d5cba2fc02a8ca3ca6d5a882b79b09c169affbfbd3bc198b5a34fd8b9c2f314df5f8aeca7bf7f71516968609a7758557d6e0440ebf2bf4c385df70eb218e7260199a5680cd333da64e4c68b3ba5fa9a1c5c816fb0bbd897a3b72e1380d8b8279d2b3e31d1e680e5609d08b7f16d\"\n}\n"),
			},
			want:    "steamcmd_bins_win32.zip.45ece1457835ba765444814e10b251645e1a8827",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.os, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
