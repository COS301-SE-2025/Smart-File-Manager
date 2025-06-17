from message_structure_pb2 import DirectoryResponse, Directory
from concurrent.futures import ThreadPoolExecutor
from metadata_scraper import MetaDataScraper
from message_structure_pb2 import DirectoryRequest, Directory, File, Tag, MetadataEntry
from kw_extractor import KWExtractor
from full_vector import FullVector
from vocabulary import Vocabulary
from k_means import KMeansCluster
import os

#temp
import numpy as np

# Master class
# Allows submission of gRPC requests. 
# Takes submitted gRPC requests and assigns them to a slave for processing before returning the response
class Master():

    def __init__(self, maxSlaves):
        self.slaves = ThreadPoolExecutor(maxSlaves)
        self.scraper = MetaDataScraper()
        self.kw_extractor = KWExtractor()
        self.vocab = Vocabulary()
        self.full_vec = FullVector()        


    # Takes gRPC request's root and sends it to be processed by a slave
    def submitTask(self, request : DirectoryRequest):
        future = self.slaves.submit(self.process, request)
        return future # Return future so non blocking

    def process(self, request : DirectoryRequest) -> DirectoryResponse:

        # Modifies directory request by adding metadata
        self.scrapeMetadata(request.root)
        # Returns top x keywords from each text file
        kw_response = self.kw_extractor.extract_kw(request.root)
        # Create a vocabulary from all keywords gathered # -> may be an issue with concurrency (might need global vocab?) -> addToVocab?
        vocabKW = self.vocab.createVocab(kw_response)
        # how will we ensure that there are always exact matches?
        full_vec = self.full_vec.createFullVector(kw_response,vocabKW, request.root)
        for vec in full_vec:
            print("\n")
            print(vec)
            print("\n")
        if len(full_vec) > 0:
            kmeans = KMeansCluster(3) # somewhat random...
            labels = kmeans.cluster(full_vec)            
            print("Labels: ", labels)
            ## temp ###############
            # Here vectors have length = X keyword + Y filetype + 1 size = X + Y + 1 features
            # Exact data passed in so it will be 100% match
            # But our "training" is our "test" aswell ? -> Ask Phil
            new_points = [
                [np.float64(0.47934089813688097), np.float64(0.20387059463771529), np.float64(0.7277683274730531), np.float64(0.8320501209552513), np.float64(0.12896877083009484), np.float64(0.9415433600539296), np.float64(0.88425354451953), np.float64(0.7698327962054629), np.float64(0.8116511751352949), np.float64(0.9548413276207754), np.float64(0.8854679418532557), np.float64(0.9783318406099909), np.float64(0.47680333045641765), np.float64(0.09035801182181236), np.float64(0.7925191634577741), np.float64(0.38306190800664985), np.float64(0.13319340059493134), np.float64(0.9531216855009049), np.float64(0.05152618173510315), np.float64(0.6226641619880533), np.float64(0.38557913607087), np.float64(0.7252901040761862), np.float64(0.8746008451139917), np.float64(0.5111429040918843), np.float64(0.2500297620776135), np.float64(0.8658982500136124), np.float64(0.12322481108802652), np.float64(0.17815978265687216), np.float64(0.95686572964321), np.float64(0.6719158031333518), np.float64(0.12426383601061686), np.float64(0.9989626075365498), np.float64(0.9263059683899479), np.float64(0.17286289632333618), np.float64(0.65134474623259), np.float64(0.1448919765993898), np.float64(0.9930327661121131), np.float64(0.4544758006067653), np.float64(0.41647516902770854), np.float64(0.6050631309366022), np.float64(0.13378373931661713), np.float64(0.9732375067694193), np.float64(0.9632163073412521), np.float64(0.07415194103630773), np.float64(0.8580549603759903), np.float64(0.24110799411440276), np.float64(0.4105081403074776), np.float64(0.43826960759452016), np.float64(0.0531805397266889), np.float64(0.8663943165491105), np.float64(0.13346680360550478), np.float64(0.727127683144778), np.float64(0.027689766176854258), np.float64(0.9953889183699364), np.float64(0.13018564642199115), np.float64(0.6026845275887566), np.float64(0.8455209950503539), np.float64(0.29162651634027503), np.float64(0.7995729353483852), np.float64(0.6347817026228165), 0.0, 1.0, 0.0, 48986] ,
                [np.float64(0.7669149614714175), np.float64(0.39570331668796077), np.float64(0.0764915994886256), np.float64(0.5406861258081728), np.float64(0.0016243052957214399), np.float64(0.9293516749510454), np.float64(0.3767539327753726), np.float64(0.9929560168327912), np.float64(0.345543913268468), np.float64(0.703553211779363), np.float64(0.2882215749884992), np.float64(0.3323413150630464), np.float64(0.17919007783821783), np.float64(0.7076097775033352), np.float64(0.018178684553718405), np.float64(0.21586543398949587), np.float64(0.26060577461100176), np.float64(0.7792050529143559), np.float64(0.8281116748106844), np.float64(0.2674120959119334), np.float64(0.8910731536102698), np.float64(0.2347921208202136), np.float64(0.45246676696379884), np.float64(0.8976674767253543), np.float64(0.490173508749583), np.float64(0.7595987075739672), np.float64(0.6646516271568645), np.float64(0.08957457555252113), np.float64(0.36012534583643907), np.float64(0.44295224590986104), np.float64(0.390248579418522), np.float64(0.7579103988809414), np.float64(0.6596273085756298), np.float64(0.728262711591656), np.float64(0.748552169378274), np.float64(0.8089186961344212), np.float64(0.13074831987240387), np.float64(0.29038153609604944), np.float64(0.4098735603127007), np.float64(0.5548998845492279), np.float64(0.5321850987162137), np.float64(0.48345533903796734), np.float64(0.8489856596412584), np.float64(0.2968705083957871), np.float64(0.017873004948633198), np.float64(0.03116082834639544), np.float64(0.8554023133302668), np.float64(0.41332386505647056), np.float64(0.4526678104777152), np.float64(0.929472304107382), np.float64(0.31582196629822856), np.float64(0.3595423109857957), np.float64(0.15273318358626475), np.float64(0.8427460496419131), np.float64(0.15291377178915455), np.float64(0.5081379281195518), np.float64(0.27888589186845403), np.float64(0.8116086888943239), np.float64(0.9681099358140003), np.float64(0.2581739924652121), 0.0, 0.0, 1.0, 90692] ,
                [np.float64(0.11957628610039883), np.float64(0.3040368320856429), np.float64(0.49514581384227585), np.float64(0.07827709699702701), np.float64(0.2783906457955726), np.float64(0.43254396189632893), np.float64(0.07289075453346783), np.float64(0.5839130759617638), np.float64(0.7710784645034068), np.float64(0.9851071125162293), np.float64(0.020410739572244374), np.float64(0.7636169960869394), np.float64(0.7508352890098873), np.float64(0.1333530840944751), np.float64(0.8062416815566585), np.float64(0.800992460814911), np.float64(0.19891115215628796), np.float64(0.2522008714856412), np.float64(0.6758244586256117), np.float64(0.4506427072443683), np.float64(0.9707821703537521), np.float64(0.33416723666512926), np.float64(0.4256900566855436), np.float64(0.09049487768294984), np.float64(0.8542521040773492), np.float64(0.839090930775108), np.float64(0.6552473796429651), np.float64(0.1083021712752974), np.float64(0.5170436524064026), np.float64(0.037887470633379716), np.float64(0.3208704071170183), np.float64(0.9732301014713941), np.float64(0.8038747786904751), np.float64(0.9158190898113687), np.float64(0.9366662709897243), np.float64(0.8602902989318818), np.float64(0.02008882833500447), np.float64(0.39843784261644255), np.float64(0.5965530522585936), np.float64(0.0011956282200019652), np.float64(0.01899109528544385), np.float64(0.3383720763621183), np.float64(0.47609753087387563), np.float64(0.8834090288942797), np.float64(0.04595314836271114), np.float64(0.14909705164987896), np.float64(0.690211672064939), np.float64(0.7553670707067989), np.float64(0.47770274382894495), np.float64(0.39156337709022504), np.float64(0.8272944254777628), np.float64(0.7301274594742572), np.float64(0.6937101650828947), np.float64(0.6412914152978031), np.float64(0.3098597622211232), np.float64(0.8850922920654665), np.float64(0.5503485735105326), np.float64(0.7858597649065298), np.float64(0.566965691871085), np.float64(0.5773581395257457), 0.0, 1.0, 0.0, 76931] ,
                [np.float64(0.5538392896394146), np.float64(0.5074977004734637), np.float64(0.9607866972568713), np.float64(0.0024548585715592486), np.float64(0.12949747936544698), np.float64(0.8303245445784541), np.float64(0.844886365097949), np.float64(0.8039386537195413), np.float64(0.6732452039148286), np.float64(0.7864702678332739), np.float64(0.1853615337268385), np.float64(0.08835815586933182), np.float64(0.8234274796944548), np.float64(0.2877176753521403), np.float64(0.7387174399607525), np.float64(0.26884363838366987), np.float64(0.9503483655527147), np.float64(0.7361374433612334), np.float64(0.023686796504951868), np.float64(0.11927358273389832), np.float64(0.735275839022648), np.float64(0.6498043425264249), np.float64(0.9177313872386725), np.float64(0.3628046592244677), np.float64(0.18209638461936872), np.float64(0.8333817476410219), np.float64(0.6938281138639958), np.float64(0.30545503164801413), np.float64(0.3071890225456648), np.float64(0.9598594395756935), np.float64(0.3618888031727354), np.float64(0.24870789980158847), np.float64(0.05656574638702405), np.float64(0.9662239995509232), np.float64(0.8941200260024731), np.float64(0.49294323099440207), np.float64(0.49441879627939755), np.float64(0.02377587055984387), np.float64(0.6325437799520959), np.float64(0.5796755426927831), np.float64(0.15981534500218564), np.float64(0.3967164334478295), np.float64(0.06934043091136499), np.float64(0.07024021027990746), np.float64(0.6357471332774927), np.float64(0.8570575510178049), np.float64(0.2523573039040037), np.float64(0.46426048620740257), np.float64(0.7166956873742228), np.float64(0.26618969241015067), np.float64(0.19292750356674548), np.float64(0.803084922655821), np.float64(0.062316659192049983), np.float64(0.017389557576718673), np.float64(0.7865834811444569), np.float64(0.24935090966039553), np.float64(0.8534076675125812), np.float64(0.2296199796611712), np.float64(0.23677802881319587), np.float64(0.4963841470384579), 0.0, 1.0, 0.0, 23565] ,
                [np.float64(0.03348041822965542), np.float64(0.17189967125070427), np.float64(0.20869989778412068), np.float64(0.39384986361389307), np.float64(0.7773453242148279), np.float64(0.17127474878588766), np.float64(0.9985054926218586), np.float64(0.49955199411585105), np.float64(0.1084937289410206), np.float64(0.48580914824796617), np.float64(0.6817065424216527), np.float64(0.71506300289635), np.float64(0.9060841333759706), np.float64(0.7998934186668148), np.float64(0.3682686487294786), np.float64(0.8521380649191089), np.float64(0.3033112809339622), np.float64(0.5307687420902462), np.float64(0.2674331371281323), np.float64(0.9529587907156133), np.float64(0.005423948750426955), np.float64(0.38824892357408447), np.float64(0.152168761431546), np.float64(0.4943323921487496), np.float64(0.44101965083948647), np.float64(0.9095446550910669), np.float64(0.8687764489747124), np.float64(0.18223939463720273), np.float64(0.056897785211135754), np.float64(0.9839484048051808), np.float64(0.8851487288548148), np.float64(0.6356341582637659), np.float64(0.13533307631556946), np.float64(0.07376912586364015), np.float64(0.5901446940277368), np.float64(0.5440764237871856), np.float64(0.677857276588707), np.float64(0.9941422268123257), np.float64(0.3070364754653957), np.float64(0.9414330439211963), np.float64(0.8364184724780309), np.float64(0.3969612350423243), np.float64(0.08016091356415478), np.float64(0.17787843679527893), np.float64(0.10911970427100792), np.float64(0.06866460783546713), np.float64(0.9529876621649984), np.float64(0.8433852690649415), np.float64(0.41205902423054097), np.float64(0.251610272391375), np.float64(0.6422721749842526), np.float64(0.8359374331665761), np.float64(0.6757548882440889), np.float64(0.012592218510113162), np.float64(0.3051990645308734), np.float64(0.22205656733064894), np.float64(0.7759636241087459), np.float64(0.8629122984427583), np.float64(0.4142628973896636), np.float64(0.13724165370351182), 0.0, 1.0, 0.0, 78693]
                ]
            predictions,centers = kmeans.predict(new_points)
            np.set_printoptions(suppress=True, precision=5) # probably temporary
            print("Pred: ", predictions)
            print("Centers: ")    
            for cent in centers:
                print(cent, "\n")
            ##################################################
        else:
            print("Invalid data")
            return response
        # response_directory = kmeans.cluster
        response_directory = request.root
        response = DirectoryResponse(root=response_directory)
        return response
    
        '''
        What is this method supposed to really do?
        * Extract metadata
        * Perform clustering
        * Generate new tree strucutre as gRPC Directory

        Preferably each of the above should be handled by its own class for seperation of concerns
        '''

    # Traverses Directory recursively and extracts metadata for each file
    def scrapeMetadata(self, currentDirectory : Directory) -> None:

        # Extract metadata
        for curFile in currentDirectory.files:
            # Ensure file path is valid
            try:
                self.scraper.set_file(os.path.abspath(curFile.original_path))
            except ValueError:
                # Invalid path => add error tag to metadata entry
                meta_error = MetadataEntry(key="Error", value="File does not exist - could not extract metadata")
                curFile.metadata.append(meta_error)
                continue
            else:
                # Valid path => scrape
                self.scraper.get_metadata()
                extracted_metadata = self.scraper.metadata
                for k,v in extracted_metadata.items():

                    meta_entry = MetadataEntry(key=str(k), value = str(v))
                    curFile.metadata.append(meta_entry)

        # Recurisve call
        if len(currentDirectory.directories) != 0:
            for curDir in currentDirectory.directories:
                self.scrapeMetadata(curDir)

######################## Temp for testing
tag1 = Tag(name="ImFixed")
meta1 = MetadataEntry(key="author", value="johnny")
meta4 = MetadataEntry(key="mime_type", value="text/plain")
meta2 = MetadataEntry(key="mime_type", value="application/pdf")
meta3 = MetadataEntry(key="mime_type", value="application/msword")

file1 = File(
    name="gopdoc.pdf",
    original_path="python/testing/test_files/myPdf.pdf",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta2]
)
file2 = File(
    name="gopdoc2.pdf",
    original_path="python/testing/test_files/testFile.txt",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta4]
)
file3 = File(
    name="gopdoc2.pdf",
    original_path="python/testing/test_files/myWordDoc.docx",
    new_path="/usr/trash/gopdoc.pdf",
    tags=[tag1],
    metadata=[meta1, meta3]
)

dir1 = Directory(
    name="useless_files",
    path="/usr/trash",
    files=[file1, file2,file3],
    directories=[]
)
req = DirectoryRequest(root=dir1) 

if __name__ == "__main__":
    master = Master(1)
    master.process(req)