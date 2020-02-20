#ifdef __cplusplus
extern "C" {
#endif
    typedef struct MetadataItem {
        char* character;
        int timestep;
        float start_time;
    } MetadataItem;

    typedef struct Metadata {
        MetadataItem* items;
        int num_items;
        double confidence;
    } Metadata;

    typedef void* ModelWrapper;
    ModelWrapper* New(const char* aModelPath, int aBeamWidth);
    void Close(ModelWrapper* w);
    void SetModelBeamWidth(ModelWrapper* w, unsigned int aBeamWidth);
    unsigned int GetModelBeamWidth(ModelWrapper* w);
    void EnableExternalScorer(ModelWrapper* w, const char* aScorerPath);
    void DisableExternalScorer(ModelWrapper* w);
    int GetModelSampleRate(ModelWrapper* w);
    char* STT(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize);
    Metadata* STTWithMetadata(ModelWrapper* w, const short* aBuffer, unsigned int aBufferSize);

    typedef void* StreamWrapper;
    StreamWrapper* CreateStream(ModelWrapper* w);
    void FreeStream(StreamWrapper* sw);
    void FeedAudioContent(StreamWrapper* sw, const short* aBuffer, unsigned int aBufferSize);
    char* IntermediateDecode(StreamWrapper* sw);
    char* FinishStream(StreamWrapper* sw);
    Metadata* FinishStreamWithMetadata(StreamWrapper* sw);

    MetadataItem* Metadata_GetItems(Metadata* m);
    double Metadata_GetConfidence(Metadata* m);
    int Metadata_GetNumItems(Metadata* m);

    char* MetadataItem_GetCharacter(MetadataItem* mi);
    int MetadataItem_GetTimestep(MetadataItem* mi);
    float MetadataItem_GetStartTime(MetadataItem* mi);

    void FreeString(char* s);
    void FreeMetadata(Metadata* m);
    char* GetVersion();

#ifdef __cplusplus
}
#endif
